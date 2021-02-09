package mysql

import (
	"fmt"

	"github.com/firadio/goutils/mysql/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Class struct {
	Db *sqlx.DB
}

func New(dsn string) Class {
	Db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("mysql connect failed, detail is [%v]", err.Error())
	}
	this := Class{}
	this.Db = Db
	return this
}

func (this *Class) Sql() sql.Class {
	sql1 := sql.New(this)
	return sql1
}

type MysqlQueryCB func([]string)

func (this *Class) Query(sql string, fun MysqlQueryCB) error {
	rows, err := this.Db.Query(sql)
	if err != nil {
		fmt.Printf("query faied, error:[%v]", err.Error())
		return err
	}
	cols, err := rows.Columns()
	if err != nil {
		fmt.Printf("rows.Columns(), error:[%v]", err.Error())
		return err
	}
	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))
	dest := make([]interface{}, len(cols))
	for i, _ := range cols {
		//fmt.Println(fieldName)
		dest[i] = &rawResult[i]
	}
	i := 0
	for rows.Next() {
		i++
		//定义变量接收查询数据
		err := rows.Scan(dest...)
		if err != nil {
			fmt.Println("get data failed, error:[%v]", err.Error())
			return err
		}
		for i, raw := range rawResult {
			if raw == nil {
				result[i] = ""
			} else {
				result[i] = string(raw)
			}
		}
		fun(result)
	}
	//关闭结果集（释放连接）
	rows.Close()
	return nil
}
