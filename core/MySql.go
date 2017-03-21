package core

import (
	"database/sql"
	"fmt"
	"regexp"
	"encoding/json"
	"time"
	"strings"
	"errors"
)

var disallowedKeywords = []string{" UNION ", " LIMIT ", ";"}

type MySql struct {
	ServiceName string `default:"mysql"`
	Host        string `default:localhost`
	Port        int `default:3360`
	Database    string
	User        string
	Password    string

	PingTimeoutMillis  int
	QueryTimeoutMillis int

	db           *sql.DB
	pingTimeout  time.Duration
	queryTimeout time.Duration
}

func (s *MySql) Kind() string {
	return "MySql"
}

func (s *MySql) Name() string {
	return s.ServiceName
}

func (s *MySql) NewСheck(req *QueryRequest) (ret Check, err error) {
	var pattern *regexp.Regexp
	if len(req.Expr) > 0 {
		pattern, err = regexp.Compile(req.Expr)
		if err != nil {
			return
		}
	}

	err = s.validateQuery(req.Query)
	if err == nil {
		query := s.limitQuery(req.Query)
		ret = mySqlCheck{info: fmt.Sprintf("q: %s, e: %s", query, req.Expr), query: query, pattern: pattern, service: s}
	}
	return
}

func (s *MySql) validateQuery(query string) error {
	var err error
	queryLowCase := strings.ToUpper(query)

	for _, keyword := range disallowedKeywords {
		if strings.Contains(queryLowCase, keyword) {
			err = errors.New(fmt.Sprintf("'%s' is disallowed for Query", keyword))
			break
		}
	}
	if !(strings.HasPrefix(queryLowCase, "SELECT ") || strings.HasPrefix(queryLowCase, "SHOW ")) {
		err = errors.New("Only SELECT/SHOW queries allowed")
	}
	return err
}

func (s *MySql) limitQuery(query string) string {
	queryLowCase := strings.ToUpper(query)
	if !(strings.HasPrefix(queryLowCase, "SELECT ")) {
		return query + " LIMIT 5"
	} else {
		return query
	}
}

func (s *MySql) Init() error {
	var err error
	if s.db == nil {
		s.db, err = sql.Open("mysql", s.buildDataSourceName())

		if err == nil {
			if s.PingTimeoutMillis > 0 {
				s.pingTimeout = time.Duration(s.PingTimeoutMillis) * time.Millisecond
				log.Debug("Ping timeout for %s is %s", s.Name(), s.pingTimeout)
			}

			if s.QueryTimeoutMillis > 0 {
				s.queryTimeout = time.Duration(s.QueryTimeoutMillis) * time.Millisecond
				log.Debug("Query timeout %s is %s", s.Name(), s.queryTimeout)
			}

			//connect
			s.ping()
		} else {
			log.Debug("Database connection of %s can't be open because of %s", err)
			s.db = nil
		}

	}
	return err
}

func (s *MySql) Close() {
	if s.db != nil {
		err := s.db.Close()
		if err != nil {
			log.Debug("Closing Database connection of %s caused error %s", s.Name, err)
		}
		s.db = nil
	} else {
		log.Debug("Database connection of %s already closed", s.Name)
	}
}

func (s *MySql) Ping() error {
	err := s.Init()
	if err == nil {
		err = s.ping()
		if err != nil {
			log.Debug("'%s' can't be reached because of %s", s.Name(), err)
		}
	}
	return err
}

func (s *MySql) ping() (error) {
	if s.pingTimeout > 0 {
		_, err := s.db.ExecContext(TimeoutContext(s.pingTimeout), "SELECT 1")
		return err
		//return s.db.PingContext(TimeoutContext(s.pingTimeout)) //not reliable, also not for GO 1.8
	} else {
		_, err := s.db.Exec("SELECT 1")
		return err
		//return s.db.Ping()
	}
}

func (s *MySql) buildDataSourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", s.User, s.Password, s.Host, s.Port, s.Database)
}

func (s *MySql) jsonBytes(sql string) (data []byte, err error) {
	obj, err := s.queryToMap(sql)
	if err == nil {
		data, err = json.Marshal(obj)
	}
	return
}

func (s *MySql) json(sql string) (json string, err error) {
	data, err := s.jsonBytes(sql)
	if err == nil {
		json = string(data)
	}
	return
}

func (s *MySql) queryToMap(sql string) ([]map[string]interface{}, error) {
	rows, err := s.query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return tableData, nil
}

func (s *MySql) query(sql string) (*sql.Rows, error) {
	if s.queryTimeout > 0 {
		return s.db.QueryContext(TimeoutContext(s.queryTimeout), sql)
	} else {
		return s.db.Query(sql)
	}
}

//check
type mySqlCheck struct {
	info    string
	query   string
	pattern *regexp.Regexp
	service *MySql
}

func (o mySqlCheck) Validate() (err error) {
	data, err := o.Query()
	if err == nil {
		if o.pattern != nil {
			if !o.pattern.Match(data) {
				err = errors.New(fmt.Sprintf("No match for %s", o.info))

			}
		} else if len(data) == 0 {
			err = errors.New(fmt.Sprintf("No match, empty result for %s", o.info))
		}
	}
	return
}

func (o mySqlCheck) Query() (data []byte, err error) {
	err = o.service.Init()
	if err != nil {
		return
	}
	data, err = o.service.jsonBytes(o.query)
	if err == nil {
		log.Debug("%s: %s", o.info, data)
	}
	return
}

func (o mySqlCheck) Close() {
}
