package dao

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dicDB = make(map[string]*DB)
)

type ScanFun func(rows *sql.Rows) interface{}
type TxFun func(tx *sql.Tx) (int64, error)
type DB struct {
	connectionStr string
	db            *sql.DB
}

func (this *DB) SetMaxIdleConns(count int) {
	this.db.SetMaxIdleConns(count)
}
func (this *DB) SetMaxOpenConns(count int) {
	this.db.SetMaxOpenConns(count)
}

func NewDB(user string, passspord string, host string) *DB {
	connectionStr := user + ":" + passspord + "@" + host
	db, ok := dicDB[connectionStr]
	if !ok {
		_db, err := sql.Open("mysql", connectionStr)
		if err != nil {
			log.Fatal(err)
		}
		_db.SetMaxIdleConns(10)
		_db.SetMaxOpenConns(10)
		db = &DB{connectionStr, _db}
		dicDB[connectionStr] = db
	}
	return db
}

func (this *DB) Select(s string, paras []interface{}, fun ScanFun) ([]interface{}, error) {
	rows, err := this.db.Query(s, paras...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	lst := make([]interface{}, 0)
	for rows.Next() {
		obj := fun(rows)
		lst = append(lst, obj)
	}
	return lst, nil
}

func (this *DB) SelectOne(s string, paras []interface{}, vals []interface{}) error {
	err := this.db.QueryRow(s, paras...).Scan(vals...)
	return err
}

func (this *DB) Insert(s string, paras ...interface{}) (int64, error) {
	stmt, err := this.db.Prepare(s)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(paras...)
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}
func (this *DB) Update(s string, paras ...interface{}) (int64, error) {
	stmt, err := this.db.Prepare(s)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(paras...)
	if err != nil {
		return 0, err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowCnt, nil
}
func (this *DB) Delete(s string, paras ...interface{}) (int64, error) {
	return this.Update(s, paras)
}
func (this *DB) Tx(fun TxFun) (int64, error) {
	tx, err := this.db.Begin()
	if err != nil {
		return 0, err
	}
	result, err := fun(tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return result, nil
}
