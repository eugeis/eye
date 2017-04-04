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
	Name string `default:"mysql"`
	AccessKey   string `default:"mysql"`

	Host string `default:localhost`
	Port int    `default:3306`

	Database string

	PingTimeoutMillis  int
	QueryTimeoutMillis int
}

type MySqlService struct {
	mysql        *MySql
	accessFinder AccessFinder

	db           *sql.DB
	pingTimeout  time.Duration
	queryTimeout time.Duration
}

func (o *MySqlService) Name() string {
	return o.mysql.Name
}

func (o *MySqlService) validateQuery(query string) error {
	var err error
	queryLowCase := strings.ToUpper(query)

	for _, keyword := range disallowedKeywords {
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

func (o *MySqlService) limitQuery(query string) string {
	queryLowCase := strings.ToUpper(query)
	if !(strings.HasPrefix(queryLowCase, "SELECT ")) {
		return query + " LIMIT 5"
	} else {
		return query
	}
}

func (o *MySqlService) Init() (err error) {
	if o.db == nil {
		var access Access
		access, err = o.accessFinder.FindAccess(o.mysql.AccessKey)
		if err == nil {
			dataSource := fmt.Sprintf("%v:%s@tcp(%v:%d)/%v", access.User, access.Password,
				o.mysql.Host, o.mysql.Port, o.mysql.Database)
			o.db, err = sql.Open("mysql", dataSource)

			if err == nil {
				if o.mysql.PingTimeoutMillis > 0 {
					o.pingTimeout = time.Duration(o.mysql.PingTimeoutMillis) * time.Millisecond
					l.Debug("Ping timeout for %v is %v", o.Name(), o.pingTimeout)
				}

				if o.mysql.QueryTimeoutMillis > 0 {
					o.queryTimeout = time.Duration(o.mysql.QueryTimeoutMillis) * time.Millisecond
					l.Debug("Query timeout %v is %v", o.Name(), o.queryTimeout)
				}

				//connect
				o.ping()
			} else {
				l.Debug("Database connection of %v can't be open because of %v", err)
				o.db = nil
			}
		}
	}
	return
}

func (o *MySqlService) Close() {
	if o.db != nil {
		err := o.db.Close()
		if err != nil {
			l.Debug("Closing Database connection of %v caused error %v", o.Name, err)
		}
		o.db = nil
	} else {
		l.Debug("Database connection of %v already closed", o.Name())
	}
}

func (o *MySqlService) Ping() error {
	err := o.Init()
	if err == nil {
		err = o.ping()
		if err != nil {
			l.Debug("'%v' can't be reached because of %v", o.Name(), err)
		}
	}
	return err
}

func (o *MySqlService) ping() error {
	return o.pingByQuery()
}

/* not reliable, also for GO 1.8? */
func (o *MySqlService) pingByConnection() error {
	if o.pingTimeout > 0 {
		return o.db.PingContext(TimeoutContext(o.pingTimeout))
	} else {
		return o.db.Ping()
	}
}

func (o *MySqlService) pingByQuery() error {
	if o.pingTimeout > 0 {
		_, err := o.db.ExecContext(TimeoutContext(o.pingTimeout), "SELECT 1")
		return err
	} else {
		_, err := o.db.Exec("SELECT 1")
		return err
	}
}

func (o *MySqlService) jsonBytes(sql string) (ret QueryResult, err error) {
	data, err := o.queryToMap(sql)
	if err == nil && len(data) > 0 {
		ret, err = json.Marshal(data)
	}
	return
}

func (o *MySqlService) json(sql string) (json string, err error) {
	data, err := o.jsonBytes(sql)
	if err == nil && len(data) > 0 {
		json = string(data)
	}
	return
}

func (o *MySqlService) queryToMap(sql string) ([]map[string]interface{}, error) {
	rows, err := o.query(sql)
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

func (o *MySqlService) query(sql string) (*sql.Rows, error) {
	if o.queryTimeout > 0 {
		return o.db.QueryContext(TimeoutContext(o.queryTimeout), sql)
	} else {
		return o.db.Query(sql)
	}
}

func (o *MySqlService) NewСheck(req *QueryRequest) (ret Check, err error) {
	var pattern *regexp.Regexp
	if len(req.Expr) > 0 {
		pattern, err = regexp.Compile(req.Expr)
		if err != nil {
			return
		}
	}

	err = o.validateQuery(req.Query)
	if err == nil {
		query := o.limitQuery(req.Query)
		ret = mySqlCheck{info: req.CheckKey(o.Name()), query: query, pattern: pattern, service: o}
	}
	return
}

//buildCheck
type mySqlCheck struct {
	info    string
	query   string
	pattern *regexp.Regexp
	service *MySqlService
}

func (o mySqlCheck) Info() string {
	return o.info
}

func (o mySqlCheck) Validate() (err error) {
	data, err := o.Query()
	if err == nil {
		if o.pattern != nil {
			if !o.pattern.Match(data) {
				err = errors.New(fmt.Sprintf("No match for %v", o.info))
			}
		} else if data == nil {
			err = errors.New(fmt.Sprintf("No match, empty result for %v", o.info))
		}
	}
	return
}

func (o mySqlCheck) Query() (data QueryResult, err error) {
	err = o.service.Init()
	if err == nil {
		data, err = o.service.jsonBytes(o.query)
		//l.Debug(string(data))
	}
	return
}
