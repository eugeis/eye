package core

import (
	"fmt"
	"strings"
	"time"
	"errors"
	"context"
	"regexp"
	"io"
	"gopkg.in/Knetic/govaluate.v2"
	"strconv"
	"encoding/json"
)

type Service interface {
	Name() string

	Init() error
	Close()

	Ping() error

	NewСheck(req *ValidationRequest) (Check, error)

	NewExporter(req *ExportRequest) (Exporter, error)
}

type Check interface {
	Info() string
	Query() (QueryResults, error)
	Validate() error
}

type Exporter interface {
	Info() string
	Export(params map[string]string) error
}

type Factory interface {
	Find(name string) (Service, error)
	Close()
}

type ValidationRequest struct {
	Query    string
	RegExpr  string
	EvalExpr string
	All      bool
}

type ExportRequest struct {
	Query     string
	Converter func(map[string]interface{}) []byte
	Out       func(params map[string]string) (io.WriteCloser, error)
}

func (o *ValidationRequest) CheckKey(serviceName string) string {
	return fmt.Sprintf("%v.q(%v).e(%v).eval(%v)",
		serviceName, o.Query, o.RegExpr, o.EvalExpr)
}

func (o *ValidationRequest) ChecksKey(strictness string, serviceNames []string) string {
	return fmt.Sprintf("%v(%v.q(%v).e(%v).eval(%v))", strictness,
		strings.Join(serviceNames, "-"), o.Query, o.RegExpr, o.EvalExpr)
}

func (o *ExportRequest) ExportKey(serviceName string) string {
	return fmt.Sprintf("%v.q(%v)", serviceName, o.Query)
}

type QueryResults []QueryResult
type QueryResult interface {
	Get(name string) (interface{}, error)
	String() string
}

func (o *QueryResults) String() (ret string) {
	if data, err := json.Marshal(o); err == nil {
		ret = string(data)
	}
	return
}

type MapQueryResult struct {
	Data map[string]interface{}
}

func (o *MapQueryResult) Get(name string) (ret interface{}, err error) {
	ret = o.Data[name]
	return
}

func (o *MapQueryResult) String() (ret string) {
	if data, err := json.Marshal(o.Data); err == nil {
		ret = string(data)
	}
	return
}

type CompositeQueryResult struct {
	Splitter *regexp.Regexp
	Data     []QueryResult
}

func NewCompositeQueryResult(separator string, data []QueryResult) (ret *CompositeQueryResult, err error) {
	var splitter *regexp.Regexp
	if splitter, err = buildSplitter(separator); err == nil {
		ret = &CompositeQueryResult{Splitter: splitter, Data: data}
	}
	return
}
func buildSplitter(separator string) (ret *regexp.Regexp, err error) {
	return regexp.Compile(fmt.Sprintf("(\\d*)%v(.*)", separator))
}

func (o *CompositeQueryResult) Get(name string) (ret interface{}, err error) {
	if indexKey := o.Splitter.Split(name, 2); len(indexKey) == 2 {
		index, _ := strconv.Atoi(indexKey[0])
		if index > 0 && index <= len(o.Data) {
			current := o.Data[index]
			ret, err = current.Get(indexKey[1])
		} else {
			err = errors.New(fmt.Sprintf("The index '%v' is not in range of '[1-%v]: ", index, len(o.Data)))
		}
	} else {
		err = errors.New("The key does not apply to the pattern: " + o.Splitter.String())
	}
	return
}

func (o *CompositeQueryResult) String() (ret string) {
	if data, err := json.Marshal(o.Data); err == nil {
		ret = string(data)
	}
	return
}

func ComposeQueryResults(separator string, results ...QueryResults) (ret QueryResults, err error) {
	var splitter *regexp.Regexp
	if splitter, err = buildSplitter(separator); err != nil {
		return
	}

	maxLen := 0
	for _, item := range results {
		if maxLen < len(item) {
			maxLen = len(item)
		}
	}

	ret = make([]QueryResult, maxLen)
	resultsCount := len(results)
	for i := 0; i < maxLen; i++ {
		data := make([]QueryResult, resultsCount)
		ret[i] = &CompositeQueryResult{Splitter: splitter, Data: data}
		for y, item := range results {
			if i < len(item) {
				data[y] = item[i]
			}
		}
	}
	return

}

func TimeoutContext(timeout time.Duration) context.Context {
	c, _ := context.WithTimeout(context.Background(), timeout)
	return c
}

type SimpleServiceFactory struct {
	services map[string]Service
}

func NewFactory() SimpleServiceFactory {
	return SimpleServiceFactory{services: make(map[string]Service)}
}

func (o *SimpleServiceFactory) Find(name string) (ret Service, err error) {
	nameLowCase := strings.ToLower(name)
	ret = o.services[nameLowCase]
	if ret == nil {
		err = errors.New(fmt.Sprintf("Can't find the service '%v'", name))
	}
	return
}

func (o *SimpleServiceFactory) Add(service Service) {
	nameLowCase := strings.ToLower(service.Name())
	o.services[nameLowCase] = service
}

func (o *SimpleServiceFactory) Close() {
	for _, item := range o.services {
		item.Close()
	}
	o.services = make(map[string]Service)
}

func compileRegExpr(pattern string) (ret *regexp.Regexp, err error) {
	if len(pattern) > 0 {
		ret, err = regexp.Compile(pattern)
	}
	return
}

func compileEval(evalExpr string) (ret *govaluate.EvaluableExpression, err error) {
	if len(evalExpr) > 0 {
		ret, err = govaluate.NewEvaluableExpression(evalExpr)
	}
	return
}

func compilePatterns(pattern ...string) (ret []*regexp.Regexp, err error) {
	ret = make([]*regexp.Regexp, len(pattern))
	for _, p := range pattern {
		if len(p) > 0 {
			ret[0], err = regexp.Compile(p)
			if err != nil {
				break
			}
		}
	}
	return
}

func ChecksKey(strictness string, serviceNames []string) string {
	return fmt.Sprintf("%v(%v)", strictness, strings.Join(serviceNames, "-"))
}

func buildQueryResult(pattern *regexp.Regexp, data []byte) (ret QueryResults, err error) {
	//Log.Debug(string(Data))
	if pattern != nil {
		if pattern != nil {
			matches := pattern.FindAllSubmatch(data, -1)
			ret = make([]QueryResult, len(matches))
			for l, match := range matches {
				entry := make(map[string]interface{})
				for i, name := range pattern.SubexpNames() {
					if i != 0 {
						entry[name] = match[i]
					}
				}
				ret[l] = &MapQueryResult{Data: entry}
			}
		}
	}
	return
}

func validate(check Check, eval *govaluate.EvaluableExpression, all bool) (err error) {
	var items QueryResults
	if items, err = check.Query(); err == nil {
		err = validateData(items, eval, all, check.Info())
	}
	return
}

func validateData(items QueryResults, eval *govaluate.EvaluableExpression, all bool, info string) (err error) {
	valid := true
	if eval != nil && len(items) > 0 {
		for _, item := range items {
			var evalResult interface{}
			println(item.String())
			if evalResult, err = eval.Eval(item); err != nil {
				valid = false
			} else if evalResult != nil {
				if boolResult, ok := evalResult.(bool); ok {
					valid = boolResult
				}
			} else {
				valid = false
			}
			if (!valid && all) || valid {
				break
			}
		}
	} else if len(items) == 0 {
		err = errors.New(fmt.Sprintf("No valid, because empty result for %v", info))
	}
	if !valid {
		err = errors.New(fmt.Sprintf("No valid for %v", info))
	}
	return
}

func prepareQuery(query string, params map[string]string) (ret string) {
	ret = query
	if params != nil && len(params) > 0 {
		for k, v := range params {
			placeHolder := fmt.Sprintf("@@%v@@", strings.ToUpper(k))
			ret = strings.Replace(ret, placeHolder, v, -1)
		}
	}
	return
}
