package core

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var disallowedKeywords = []string{" UNION ", " LIMIT ", ";"}

type MySql struct {
	ServiceName string `default:"mysql"`
	AccessKey   string `default:"mysql"`

	Host string `default:localhost`
	Port int    `default:3306`

	Database string

	PingTimeoutMillis  int
	QueryTimeoutMillis int

	db           *sql.DB
	pingTimeout  time.Duration
	queryTimeout time.Duration

	accessFinder AccessFinder
}

func (s *MySql) Name() string {
	return s.ServiceName
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

func (s *MySql) Init() (err error) {
	if s.db == nil {
		var access Access
		access, err = s.accessFinder.FindAccess(s.AccessKey)
		if err == nil {
			dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", access.User, access.Password, s.Host, s.Port, s.Database)
			s.db, err = sql.Open("mysql", dataSource)

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
	}
	return
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

func (s *MySql) ping() error {
	return s.pingByQuery()
}

/* not reliable, also for GO 1.8? */
func (s *MySql) pingByConnection() error {
	if s.pingTimeout > 0 {
		return s.db.PingContext(TimeoutContext(s.pingTimeout))
	} else {
		return s.db.Ping()
	}
}

func (s *MySql) pingByQuery() error {
	if s.pingTimeout > 0 {
		_, err := s.db.ExecContext(TimeoutContext(s.pingTimeout), "SELECT 1")
		return err
	} else {
		_, err := s.db.Exec("SELECT 1")
		return err
	}
}

func (s *MySql) jsonBytes(sql string) (ret QueryResult, err error) {
	data, err := s.queryToMap(sql)
	if err == nil && len(data) > 0 {
		ret, err = json.Marshal(data)
	}
	return
}

func (s *MySql) json(sql string) (json string, err error) {
	data, err := s.jsonBytes(sql)
	if err == nil && len(data) > 0 {
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
		ret = mySqlCheck{info: req.CheckKey(s.Name()), query: query, pattern: pattern, service: s}
	}
	return
}

//buildCheck
type mySqlCheck struct {
	info    string
	query   string
	pattern *regexp.Regexp
	service *MySql
}

func (o mySqlCheck) Info() string {
	return o.info
}

func (o mySqlCheck) Validate() (err error) {
	data, err := o.Query()
	if err == nil {
		if o.pattern != nil {
			if !o.pattern.Match(data) {
				err = errors.New(fmt.Sprintf("No match for %s", o.info))
			}
		} else if data == nil {
			err = errors.New(fmt.Sprintf("No match, empty result for %s", o.info))
		}
	}
	return
}

func (o mySqlCheck) Query() (data QueryResult, err error) {
	err = o.service.Init()
	if err == nil {
		data, err = o.service.jsonBytes(o.query)
		//log.Debug(string(data))
	}
	return
}
