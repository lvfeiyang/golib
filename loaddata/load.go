package loaddata

import (
	"database/sql"
	"fmt"
	"time"
)

type DataTo interface {
	HandleOneLine([]interface{}) error
	LoadEnd()
}

func Loadinit(dbci interface{}, writeTo DataTo) error {
	go func() {
		c := time.Tick(5 * time.Minute)
		for range c {
			loadData(dbci, writeTo)
		}
	}()
	return loadData(dbci, writeTo)
}

func loadData(dbci interface{}, writeTo DataTo) error {
	defer writeTo.LoadEnd()
	switch dbc := dbci.(type) {
	case *MysqlC:
		return loadFromMysql(dbc, writeTo)
	default:
		return fmt.Errorf("unknow dbc type %T", dbc)
	}
}

type MysqlC struct {
	Ip   string
	Port int

	User, Password, DB string

	Sql string
}

//文件加载 从mysql file 到 [][]string
func loadFromMysql(dbc *MysqlC, writeTo DataTo) error {
	recDb, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&timeout=30s",
		dbc.User, dbc.Password, dbc.Ip, dbc.Port, dbc.DB))
	if err != nil {
		return err
	}
	defer recDb.Close()

	rows, err := recDb.Query(dbc.Sql) //`select * from content where status = 1;`
	if err != nil {
		return err
	}
	defer rows.Close()

	var oneLine []interface{}
	var errCount int
	for rows.Next() {
		if err := rows.Scan(oneLine...); err == nil {
			if err := writeTo.HandleOneLine(oneLine); err != nil {
				errCount++
			}
		} else {
			errCount++
		}
	}
	if errCount > 0 {
		return fmt.Errorf("err lines num: %d", errCount)
	}
	return nil
}
