package test

import (
	"database/sql"
	"restapi/main/dao/model"
	"restapi/mysql/dao"
	"time"
)

var (
	db *dao.DB
)

func Init(user, password, host string) {
	db = dao.NewDB(user, password, host)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
}

func SelectTest(id int64) (*model.Test, error) {
	lst, err := db.Select("select id,`key`,`value`,update_date,create_date from tests where id=?", []interface{}{id},
		func(rows *sql.Rows) interface{} {
			t := &model.Test{}
			var u string
			var c string
			err := rows.Scan(&t.ID, &t.Key, &t.Value, &u, &c)
			if err != nil {
				return nil
			}
			t.UpdateDate = dao.GetDateTime(u, "")
			t.CreateDate = dao.GetDateTime(c, "")
			return t
		})
	if err != nil {
		return nil, err
	}
	return lst[0].(*model.Test), nil
}
func SelectByID(id int64) (*model.Test, error) {
	t := &model.Test{}
	var u string
	var c string
	err := db.SelectOne("select id,`key`,`value`,update_date,create_date from tests where id=?", []interface{}{id}, []interface{}{&t.ID, &t.Key, &t.Value, &u, &c})
	if err != nil {
		return nil, err
	}
	t.UpdateDate = dao.GetDateTime(u, "")
	t.CreateDate = dao.GetDateTime(c, "")
	return t, nil
}
func SelectByKey(key string) (*model.Test, error) {
	t := &model.Test{}
	var u string
	var c string
	err := db.SelectOne("select id,`key`,`value`,update_date,create_date from tests where `key`=?", []interface{}{key}, []interface{}{&t.ID, &t.Key, &t.Value, &u, &c})
	if err != nil {
		return nil, err
	}
	t.UpdateDate = dao.GetDateTime(u, "")
	t.CreateDate = dao.GetDateTime(c, "")
	return t, nil
}

func Insert(key string, value string) (*model.Test, error) {
	t := time.Now()
	id, err := db.Insert("insert into tests(`key`,`value`,update_date,create_date) values(?,?,?,?)", key, value, t, t)
	if err == nil {
		return &model.Test{ID: id, Key: key, Value: value, UpdateDate: t, CreateDate: t}, nil
	}
	return nil, err
}

func UpdateValue(id int64, value string) (int64, error) {
	id, err := db.Update("update tests set `value`=? where id=?", value, id)
	return id, err
}
func UpdateKey(id int64, key string) (int64, error) {
	id, err := db.Update("update tests set `key`=? where id=?", key, id)
	return id, err
}
