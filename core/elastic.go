package core

import (
	"gopkg.in/olivere/elastic.v5"
	"gee/as"
	"strings"
	"fmt"
	"regexp"
	"golang.org/x/net/context"
	"time"
	"errors"
)

var disallowedElasticKeywords = []string{}

type Elastic struct {
	Name       string `default:"elastic"`
	AccessKey  string `default:"elastic"`
	Host       string `default:"localhost"`
	Port       int    `default:"9300"`
	ScrollSize int `default:"500"`

	Index string

	PingTimeoutMillis  int
	QueryTimeoutMillis int
}

type ElasticService struct {
	elastic      *Elastic
	accessFinder as.AccessFinder

	client        *elastic.Client
	clusterHealth *elastic.ClusterHealthService
	search        *elastic.SearchService
	scroll        *elastic.ScrollService
	context       context.Context

	pingTimeout  time.Duration
	queryTimeout time.Duration
}

func (o *ElasticService) Name() string {
	return o.elastic.Name
}

func (o *ElasticService) validateQuery(query string) error {
	var err error
	queryLowCase := strings.ToUpper(query)

	for _, keyword := range disallowedElasticKeywords {
		if strings.Contains(queryLowCase, keyword) {
			err = errors.New(fmt.Sprintf("'%v' is disallowed for Query", keyword))
			break
		}
	}
	if !(strings.HasPrefix(queryLowCase, "SELECT ") || strings.HasPrefix(queryLowCase, "SHOW ")) {
		err = errors.New("Only SELECT/SHOW queries allowed")
	}
	return err
}

func (o *ElasticService) limitQuery(query string) string {
	queryLowCase := strings.ToUpper(query)
	if !(strings.HasPrefix(queryLowCase, "SELECT ")) {
		return query + " LIMIT 5"
	} else {
		return query
	}
}

func (o *ElasticService) Init() (err error) {
	if o.client == nil {
		url := fmt.Sprintf("http://(%v:%d)", o.elastic.Host, o.elastic.Port)
		o.client, err = elastic.NewClient(elastic.SetURL(url))
		if err == nil {
			o.clusterHealth = o.client.ClusterHealth().Index(o.elastic.Index)
			o.search = o.client.Search(o.elastic.Index)
			o.scroll = o.client.Scroll(o.elastic.Index).Size(o.elastic.ScrollSize)
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
			Log.Debug("Index connection of %v can't be open because of %v", err)
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
		o.search = nil
		o.scroll = nil
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

func (o *ElasticService) ping() error {
	return o.clusterHealth.Validate()
}

func (o *ElasticService) NewСheck(req *QueryRequest) (ret Check, err error) {
	var pattern *regexp.Regexp
	if pattern, err = compilePattern(req.Expr); err == nil {
		if err = o.validateQuery(req.Query); err == nil {
			query := o.limitQuery(req.Query)
			ret = elasticCheck{info: req.CheckKey(o.Name()), query: query,
				pattern:         pattern, service: o, not: req.Not}
		}
	}
	return
}

//buildCheck
type elasticCheck struct {
	info    string
	query   string
	not     bool
	pattern *regexp.Regexp
	service *ElasticService
}

func (o elasticCheck) Info() string {
	return o.info
}

func (o elasticCheck) Validate() error {
	return validate(o, o.pattern, o.not)
}

func (o elasticCheck) Query() (data QueryResult, err error) {
	if err = o.service.Init(); err == nil {
		//l.Debug(string(data))
	}
	return
}
