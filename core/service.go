package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/Knetic/govaluate.v2"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	Name() string

	Init() error
	Close()

	Ping() error

	NewСheck(req *ValidationRequest) (Check, error)

	NewExporter(req *ExportRequest) (Exporter, error)
}

type Query interface {
	Info() string
	Query() (QueryResults, error)
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

type MapWriter interface {
	WriteMap(data map[string]interface{}) error
}

type QueryResultMapWriter struct {
	Data []QueryResult
}

func NewQueryResultMapWriter() *QueryResultMapWriter {
	return &QueryResultMapWriter{Data: make([]QueryResult, 0)}
}

func (o *QueryResultMapWriter) WriteMap(data map[string]interface{}) error {
	o.Data = append(o.Data, &MapQueryResult{data})
	return nil
}

type WriteCloserMapWriter struct {
	Convert func(map[string]interface{}) io.Reader
	Out     io.WriteCloser
}

func (o *WriteCloserMapWriter) WriteMap(data map[string]interface{}) (err error) {
	_, err = io.Copy(o.Out, o.Convert(data))
	return
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

func NewValidationRequest(query string, evalExp string) *ValidationRequest {
	return &ValidationRequest{Query: query, EvalExpr: evalExp, All: true}
}

type ExportRequest struct {
	Query     string
	Convert   func(map[string]interface{}) io.Reader
	CreateOut func(params map[string]string) (io.WriteCloser, error)
}

func (o *ValidationRequest) CheckKey(serviceName string) string {
	if o.All {
		return fmt.Sprintf("%v.q(%v).e(%v).all[eval(%v)]",
			serviceName, o.Query, o.RegExpr, o.EvalExpr)
	} else {
		return fmt.Sprintf("%v.q(%v).e(%v).any[eval(%v)]",
			serviceName, o.Query, o.RegExpr, o.EvalExpr)
	}
}

func (o *ValidationRequest) ChecksKey(strictness string, serviceNames []string) string {
	if o.All {
		return fmt.Sprintf("%v(%v.q(%v).e(%v).eval(%v))", strictness,
			strings.Join(serviceNames, "-"), o.Query, o.RegExpr, o.EvalExpr)
	} else {
		return fmt.Sprintf("%v(%v.q(%v).e(%v).eval(%v))", strictness,
			strings.Join(serviceNames, "-"), o.Query, o.RegExpr, o.EvalExpr)
	}

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
	return regexp.Compile(fmt.Sprintf("(.*)%v(\\d*)", separator))
}

func (o *CompositeQueryResult) Get(name string) (ret interface{}, err error) {
	if key_index := o.Splitter.FindStringSubmatch(name); len(key_index) == 3 {
		index, _ := strconv.Atoi(key_index[2])
		if index > 0 && index <= len(o.Data) {
			current := o.Data[index-1]
			ret, err = current.Get(key_index[1])
		} else {
			err = errors.New(fmt.Sprintf("The composite key '%v'is not compatible with the available data, "+
				"the index '%v' is not in range of '[1-%v]' composite ", name, index, len(o.Data)))
		}
	} else {
		err = errors.New("The composite key does not apply to the pattern: " + o.Splitter.String())
	}
	return
}

func (o *CompositeQueryResult) String() (ret string) {
	if data, err := json.Marshal(o.Data); err == nil {
		ret = string(data)
	}
	return
}

func ComposeQueryResults(separator string, results []QueryResults) (ret QueryResults, err error) {
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

func NewFactory() *SimpleServiceFactory {
	return &SimpleServiceFactory{services: make(map[string]Service)}
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
	o.add(service.Name(), service)
}

func (o *SimpleServiceFactory) add(key string, service Service) {
	keyLowCase := strings.ToLower(key)
	o.services[keyLowCase] = service
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
		if ret, err = govaluate.NewEvaluableExpression(evalExpr); err != nil {
			Log.Err("The compilation of evaluable expression '%v' failed because of '%'", evalExpr, err)
		}
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
			if evalResult, err = eval.Eval(item); err != nil {
				valid = false
			} else if evalResult != nil {
				if boolResult, ok := evalResult.(bool); ok {
					valid = boolResult
				}
			} else {
				valid = false
			}
			if !valid && all {
				err = errors.New(fmt.Sprintf("Validation of '%v' failed agains %v", item, info))
				break
			} else if valid && !all {
				break
			}
		}
		if err == nil && !valid {
			err = errors.New(fmt.Sprintf("Validation of '%v' failed agains %v", items, info))
		}
	} else if len(items) == 0 {
		err = errors.New(fmt.Sprintf("No valid, because empty result for %v", info))
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
