package dao

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"fmt"
	"time"
)

func GetDateTime(str string, location string) time.Time {
	const Layout = "2006-01-02 15:04:05"
	if location == "" {
		location = "Asia/Shanghai"
	}
	loc, _ := time.LoadLocation(location)
	t, _ := time.ParseInLocation(Layout, str, loc)
	return t
}

func GetInt64(bs sql.RawBytes) int64 {
	buf := bytes.NewBuffer(bs)
	var i64 int64
	ok := binary.Read(buf, binary.BigEndian, &i64)
	if ok != nil {
		fmt.Println(ok)
	}
	return i64
}
