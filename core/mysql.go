package core

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"gee/as"
	"gopkg.in/Knetic/govaluate.v2"
	"strconv"
)

var disallowedSqlKeywords = []string{" UNION ", " LIMIT ", ";"}

type MySql struct {
	Name      string `default:"mysql"`
	AccessKey string `default:"mysql"`

	Host string `default:"localhost"`
	Port int    `default:"3306"`

	Database string

	PingTimeoutMillis  int
	QueryTimeoutMillis int
}

type MySqlService struct {
	mysql        *MySql
	accessFinder as.AccessFinder

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

	for _, keyword := range disallowedSqlKeywords {
		if strings.Contains(queryLowCase, keyword) {
			err = errors.New(fmt.Sprintf("'%v' is disallowed for Query", keyword))
			break
		}
	}
	if len(queryLowCase) == 0 {
		err = errors.New("Query is empOnly SELECT/SHOW queries allowed")
	}
	if !(strings.HasPrefix(queryLowCase, "SELECT ") || strings.HasPrefix(queryLowCase, "SHOW ")) {
		err = errors.New("Only SELECT/SHOW queries allowed")
	}
	return err
}

func (o *MySqlService) limitQuery(query string) string {
	queryLowCase := strings.ToUpper(query)
	if strings.HasPrefix(queryLowCase, "SELECT ") {
		return query + " LIMIT 5"
	} else {
		return query
	}
}

func (o *MySqlService) Init() (err error) {
	if o.db == nil {
		var access as.Access
		if access, err = o.accessFinder.FindAccess(o.mysql.AccessKey); err == nil {
			dataSource := fmt.Sprintf("%v:%s@tcp(%v:%d)/%v", access.User, access.Password,
				o.mysql.Host, o.mysql.Port, o.mysql.Database)
			o.db, err = sql.Open("mysql", dataSource)

			if err == nil {
				if o.mysql.PingTimeoutMillis > 0 {
					o.pingTimeout = time.Duration(o.mysql.PingTimeoutMillis) * time.Millisecond
					Log.Debug("Ping timeout for %v is %v", o.Name(), o.pingTimeout)
				}

				if o.mysql.QueryTimeoutMillis > 0 {
					o.queryTimeout = time.Duration(o.mysql.QueryTimeoutMillis) * time.Millisecond
					Log.Debug("Query timeout %v is %v", o.Name(), o.queryTimeout)
				}

				//connect
				o.ping()
			} else {
				Log.Debug("Index connection of %v can't be open because of %v", err)
				o.db = nil
			}
		}
	}
	return
}

func (o *MySqlService) Close() {
	if o.db != nil {
		if err := o.db.Close(); err != nil {
			Log.Debug("Closing Index connection of %v caused error %v", o.Name, err)
		}
		o.db = nil
	} else {
		Log.Debug("Index connection of %v already closed", o.Name())
	}
}

func (o *MySqlService) Ping() (err error) {
	if err = o.Init(); err == nil {
		if err = o.ping(); err != nil {
			Log.Debug("'%v' can't be reached because of %v", o.Name(), err)
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

/*
func (o *MySqlService) jsonBytes(sql string) (ret QueryResult, err error) {
	var Data []map[string]interface{}
	if Data, err = o.queryToMap(sql); err == nil && len(Data) > 0 {
		ret, err = json.Marshal(Data)
	}
	return
}

func (o *MySqlService) json(sql string) (json string, err error) {
	var Data QueryResult
	if Data, err = o.jsonBytes(sql); err == nil && len(Data) > 0 {
		json = string(Data)
	}
	return
}*/

func (o *MySqlService) queryToMap(sql string) (ret QueryResults, err error) {
	rows, err := o.query(sql)
	if err != nil {
		return
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return
	}
	count := len(columns)
	ret = make([]QueryResult, 0)
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
				str := string(b)
				if d, conErr := strconv.Atoi(str); conErr == nil {
					v = d
				} else {
					v = str
				}
			} else {
				v = val
			}
			entry[col] = v
		}
		ret = append(ret, &MapQueryResult{entry})
	}
	return
}

func (o *MySqlService) query(sql string) (*sql.Rows, error) {
	if o.queryTimeout > 0 {
		return o.db.QueryContext(TimeoutContext(o.queryTimeout), sql)
	} else {
		return o.db.Query(sql)
	}
}

func (o *MySqlService) NewСheck(req *ValidationRequest) (ret Check, err error) {
	var eval *govaluate.EvaluableExpression
	if eval, err = compileEval(req.EvalExpr); err != nil {
		return
	}

	if err = o.validateQuery(req.Query); err != nil {
		return
	}

	query := o.limitQuery(req.Query)
	ret = &mySqlCheck{
		info: req.CheckKey(o.Name()), query: query, service: o,
		eval: eval, all: req.All}
	return
}

func (o *MySqlService) NewExporter(req *ExportRequest) (ret Exporter, err error) {
	return
}

//buildCheck
type mySqlCheck struct {
	info    string
	query   string
	all     bool
	eval    *govaluate.EvaluableExpression
	service *MySqlService
}

func (o mySqlCheck) Info() string {
	return o.info
}

func (o mySqlCheck) Validate() error {
	return validate(o, o.eval, o.all)
}

func (o mySqlCheck) Query() (data QueryResults, err error) {
	if err = o.service.Init(); err == nil {
		data, err = o.service.queryToMap(o.query)
	}
	return
}
