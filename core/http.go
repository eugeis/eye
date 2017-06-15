package core

import (
	"errors"
	"eye/digest"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
	"gee/as"
	"gopkg.in/Knetic/govaluate.v2"
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

func (o httpCheck) Info() string {
	return o.info
}

func (o httpCheck) Validate() error {
	return validate(o, o.eval, o.all)
}

func (o httpCheck) Query() (ret QueryResults, err error) {
	err = o.service.Init()
	if err != nil {
		return
	}
	resp, err := o.req.Execute(o.service.client)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var data []byte
	if data, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	ret, err = buildQueryResult(o.pattern, data)
	Log.Debug("http ret: %s", ret)

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Status %d", resp.StatusCode))
	}
	return
}
