package core

import (
	"eye/digest"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
	"gee/as"
	"gopkg.in/Knetic/govaluate.v2"
	"io"
)

type Http struct {
	Name      string
	AccessKey string
	Url       string

	PingRequest *ValidationRequest

	PingTimeoutMillis  int
	QueryTimeoutMillis int
}

type HttpService struct {
	http         *Http
	accessFinder as.AccessFinder

	client    *http.Client
	pingCheck *httpCheck

	pingTimeout  time.Duration
	queryTimeout time.Duration
}

func (o *HttpService) Name() string {
	return o.http.Name
}

func (o *HttpService) Init() (err error) {
	if o.client == nil {
		o.client = digest.NewClient(true, o.queryTimeout)
		o.pingCheck, err = o.newСheck(o.http.PingRequest)
		if err != nil {
			o.Close()
		}
	}
	return
}

func (o *HttpService) Close() {
	o.client = nil
	o.pingCheck = nil
}

func (o *HttpService) Ping() error {
	err := o.Init()
	if err == nil {
		err = o.pingCheck.Validate()
		if err != nil {
			Log.Debug("'%v' can't be reached because of %v", o.Name(), err)
		}
	}
	return err
}

func (o *HttpService) query(req *digest.Request) (ret []byte, err error) {
	err = o.Init()
	if err != nil {
		return
	}
	resp, err := req.Execute(o.client)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	ret, err = ioutil.ReadAll(resp.Body)
	return
}

func (o *HttpService) queryToWriter(req *digest.Request, pattern *regexp.Regexp, writer MapWriter) (err error) {
	if err = o.Init(); err != nil {
		return
	}
	var resp *http.Response
	if resp, err = req.Execute(o.client); err != nil {
		return
	}
	defer resp.Body.Close()

	var data []byte
	if data, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	//Log.Debug(string(Data))
	if pattern != nil {
		matches := pattern.FindAllSubmatch(data, -1)
		for _, match := range matches {
			entry := make(map[string]interface{})
			for i, name := range pattern.SubexpNames() {
				if i != 0 {
					entry[name] = match[i]
				}
			}
			writer.WriteMap(entry)
		}
	} else {
		entry := make(map[string]interface{})
		entry["data"] = string(data)
		writer.WriteMap(entry)
	}
	return
}

func body(resp *http.Response) string {
	body, _ := ioutil.ReadAll(resp.Body)
	ret := fmt.Sprintf("%v", body)
	return ret
}

func (o *HttpService) NewСheck(req *ValidationRequest) (ret Check, err error) {
	return o.newСheck(req)
}

func (o *HttpService) newСheck(req *ValidationRequest) (ret *httpCheck, err error) {
	var access as.Access
	if access, err = o.accessFinder.FindAccess(o.http.AccessKey); err != nil {
		return
	}

	var eval *govaluate.EvaluableExpression
	if eval, err = compileEval(req.EvalExpr); err != nil {
		return
	}

	var pattern *regexp.Regexp
	if pattern, err = compileRegExpr(req.RegExpr); err != nil {
		return
	}

	dReq := digest.NewRequest(access.User, access.Password, "GET", o.http.Url+req.Query, "")
	ret = &httpCheck{
		info:    req.CheckKey(o.Name()), req: &dReq,
		pattern: pattern, service: o, eval: eval, all: req.All}
	return
}

func (o *HttpService) NewExporter(req *ExportRequest) (ret Exporter, err error) {
	ret = &httpExporter{info: req.ExportKey(o.Name()), req: req, service: o}
	return
}

//buildCheck

type httpCheck struct {
	info    string
	req     *digest.Request
	all     bool
	eval    *govaluate.EvaluableExpression
	pattern *regexp.Regexp
	service *HttpService
}

func (o *httpCheck) Info() string {
	return o.info
}

func (o *httpCheck) Validate() error {
	return validate(o, o.eval, o.all)
}

func (o *httpCheck) Query() (ret QueryResults, err error) {
	writer := NewQueryResultMapWriter()
	if err = o.service.queryToWriter(o.req, o.pattern, writer); err == nil {
		ret = writer.Data
	}
	return
}

type httpExporter struct {
	info    string
	httpReq *digest.Request
	req     *ExportRequest
	pattern *regexp.Regexp
	service *HttpService
}

func (o *httpExporter) Info() string {
	return o.info
}

func (o *httpExporter) Export(params map[string]string) (err error) {
	if err = o.service.Init(); err != nil {
		return
	}

	var out io.WriteCloser
	if out, err = o.req.CreateOut(params); err != nil {
		return
	}
	defer out.Close()

	err = o.service.queryToWriter(o.httpReq, o.pattern, &WriteCloserMapWriter{Convert: o.req.Convert, Out: out})
	return
}
