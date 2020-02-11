package golib

import (
	"github.com/rcrowley/go-bson"
	"go-android-lib/anlib"
	"runtime"
)

var goAppXid string
var GoApp GoApplication = nil

type GoApplication interface {
	AppPath() string
	SdAppPath() string
	SdCardPath() string
}

func InitGoApp(app GoApplication) {
	GoApp = app
	anlib.LogInfo("ruisgo-os", "os:", runtime.GOOS)

	initDb()

	appxid := DataGetParam("appxid")
	if appxid == nil {
		DataSetParam("appxid", []byte(bson.NewObjectId().Hex()), "string", "设备随机ID")
	}
	appxid = DataGetParam("appxid")
	goAppXid = string(appxid.Value)

	initParam()
}

func GetGoAppXid() string {
	return goAppXid
}

func initParam() {

}
