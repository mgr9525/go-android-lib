package golib

import "time"

const (
	TYPE_PARAM_DEFAULT = "default"
)

type DataParam struct {
	Id    int64  `xorm:"pk autoincr INTEGER NOT NULL"`
	Types string `xorm:"VARCHAR(50)"`
	Name  string `xorm:"unique VARCHAR(100)"`
	Title string `xorm:"VARCHAR(100)"`
	Value []byte `xorm:"BLOB"`
	Times time.Time
}
