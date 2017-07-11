package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"gee/as"
	"golang.org/x/net/context"
	"gopkg.in/Knetic/govaluate.v2"
	"gopkg.in/olivere/elastic.v5"
	"io"
	"strings"
	"time"
)

type Elastic struct {
	Name       string `default:"elastic"`
	Host       string `default:"localhost"`
	Port       int    `default:"9200"`
	ScrollSize int    `default:"500"`

	Index string

	PingTimeoutMillis  int
	QueryTimeoutMillis int
}

type ElasticService struct {
	elastic      *Elastic
	accessFinder as.AccessFinder

	client        *elastic.Client
	clusterHealth *elastic.ClusterHealthService
	context       context.Context

	pingTimeout  time.Duration
	queryTimeout time.Duration
}

func (o *ElasticService) Name() string {
	return o.elastic.Name
}

func (o *ElasticService) Init() (err error) {
	if o.client == nil {
		url := fmt.Sprintf("http://%v:%d", o.elastic.Host, o.elastic.Port)
		o.client, err = elastic.NewClient(elastic.SetURL(url))
		if err == nil {
			o.clusterHealth = o.client.ClusterHealth().Index(o.elastic.Index)
			o.context = context.Background()

			if o.elastic.PingTimeoutMillis > 0 {
				o.pingTimeout = time.Duration(o.elastic.PingTimeoutMillis) * time.Millisecond
				Log.Debug("Ping timeout for %v is %v", o.Name(), o.pingTimeout)
			}

			if o.elastic.QueryTimeoutMillis > 0 {
				o.queryTimeout = time.Duration(o.elastic.QueryTimeoutMillis) * time.Millisecond
				Log.Debug("Query timeout %v is %v", o.Name(), o.queryTimeout)
			}

			//connect
			o.ping()
		} else {
			Log.Debug("Elastic client can't connect to %v because of %v", url, err)
			o.client = nil
		}
	}
	return
}

func (o *ElasticService) Close() {
	if o.client != nil {
		o.client.Stop()
		o.clusterHealth = nil
		o.client = nil
		o.context = nil
	} else {
		Log.Debug("Client connection of %v already closed", o.Name())
	}
}

func (o *ElasticService) Ping() (err error) {
	if err = o.Init(); err == nil {
		if err = o.ping(); err != nil {
			Log.Debug("'%v' can't be reached because of %v", o.Name(), err)
		}
	}
	return err
}

func (o *ElasticService) ping() (err error) {
	var res *elastic.ClusterHealthResponse
	if res, err = o.clusterHealth.Do(o.context); err == nil {
		if strings.EqualFold(res.Status, "RED") {
			err = errors.New("Cluster health status is [RED]")
		}
	}
	return
}

func (o *ElasticService) NewExecutor(req *CommandRequest) (ret Executor, err error) {
	return nil, errors.New(fmt.Sprintf("Not implemented yet in %v", o.Name()))
}

func (o *ElasticService) NewСheck(req *ValidationRequest) (ret Check, err error) {
	var eval *govaluate.EvaluableExpression
	if eval, err = compileEval(req.EvalExpr); err != nil {
		return
	}

	ret = &elasticCheck{info: req.CheckKey(o.Name()), query: req.Query, eval: eval, all: req.All, service: o}
	return
}

func (o *ElasticService) NewExporter(req *ExportRequest) (ret Exporter, err error) {
	ret = elasticExporter{info: req.ExportKey(o.Name()), req: req, service: o}
	return
}

//buildCheck
type elasticCheck struct {
	info    string
	query   string
	all     bool
	eval    *govaluate.EvaluableExpression
	service *ElasticService
	search  *elastic.SearchService
}

func (o elasticCheck) Info() string {
	return o.info
}

func (o elasticCheck) Validate() error {
	return validate(o, o.eval, o.all)
}

func (o elasticCheck) Query() (data QueryResults, err error) {
	if err = o.service.Init(); err == nil {
		if o.search == nil {
			o.search = o.service.client.Search(o.service.elastic.Index).Size(5).Source(o.query)
		}
		//l.Debug(string(Data))
		var res *elastic.SearchResult
		if res, err = o.search.Do(o.service.context); err != nil {
			return
		}
		data := make([]QueryResult, 0)
		for _, hit := range res.Hits.Hits {
			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err == nil {
				data = append(data, &MapQueryResult{item})
			} else {
				Log.Info("Error %v, at unmarshal of %v", err, hit)
			}
		}
		Log.Debug("elastic Data: %s", data)
	}
	return
}

type elasticExporter struct {
	info    string
	req     *ExportRequest
	service *ElasticService
	scroll  *elastic.ScrollService
}

func (o elasticExporter) Info() string {
	return o.info
}

func (o elasticExporter) Export(params map[string]string) (err error) {
	if err = o.service.Init(); err != nil {
		return
	}
	if o.scroll == nil {
		o.scroll = o.service.client.Scroll(o.service.elastic.Index).Size(o.service.elastic.ScrollSize)
	}

	var out io.WriteCloser
	if out, err = o.req.CreateOut(params); err != nil {
		return
	}
	defer out.Close()

	var res *elastic.SearchResult

	query := prepareQuery(o.req.Query, params)

	o.scroll.Body(query)

	convert := o.req.Convert

	for {
		if res, err = o.search(); err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}

		var reader io.Reader
		for _, hit := range res.Hits.Hits {
			item := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &item)
			if err == nil {
				if reader, err = convert(item); err == nil {
					io.Copy(out, reader)
				} else {
					Log.Err("Can't convert item to reader because of %v", err)
				}
			} else {
				Log.Err("Can't unmarshal because of %v", err)
			}
		}
	}
	o.scroll.Clear(o.service.context)
	return
}

func (o *elasticExporter) search() (ret *elastic.SearchResult, err error) {
	ret, err = o.scroll.Do(o.service.context)
	if err == nil {
		if ret == nil {
			err = errors.New("expected results != nil; got nil")
		}
		if ret.Hits == nil {
			err = errors.New("expected results != nil; got nil")
		}
	}
	return
}
