package model

import "time"

type Test struct {
	ID         int64     `json:"id"`
	Key        string    `json:"key"`
	Value      string    `json:"value"`
	UpdateDate time.Time `json:"update_date"`
	CreateDate time.Time `json:"create_date"`
}
