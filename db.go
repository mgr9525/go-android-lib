package golib

import (
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	ruisUtil "github.com/mgr9525/go-ruisutil"
	"github.com/thinkoner/openssl"
	"github.com/xormplus/xorm"
	"go-android-lib/anlib"
	"go-android-lib/core"
	"time"
)

var dbMain *xorm.Engine

func initDb() {
	defer ruisUtil.Recovers("initDb", func(errs string) {
		anlib.LogError("ruisgo-initDb", "err:"+errs)
	})
	path := fmt.Sprintf("%s/files/main.db", GoApp.AppPath())
	db, err := xorm.NewEngine("sqlite3", path)
	if err != nil {
		anlib.LogError("ruisgo-initDb", "open db err(", path, "):", err.Error())
		return
	}
	db.Sync2(DataParam{})
	db.Sync2(DataUser{})
	dbMain = db
}
func DataGetParam(nm string) *DataParam {
	defer ruisUtil.Recovers("DataGetParam", func(errs string) {
		anlib.LogError("ruisgo-DataGetParam", "err:"+errs)
	})
	e := &DataParam{}
	ok, err := dbMain.Where("`name`=?", nm).Get(e)
	if err != nil {
		anlib.LogError("ruisgo-DataGetParam", "get err:"+err.Error())
		return nil
	}
	if ok {
		btaes, err := openssl.AesECBDecrypt(e.Value, core.DbAESKEY, openssl.PKCS7_PADDING)
		//btaes,err:=ruisUtil.AESDecrypt(e.Value,core.DbAESKEY,core.DbAESIV)
		if err != nil {
			return nil
		}
		e.Value = btaes
		return e
	}
	return nil
}
func DataSetParam(nm string, bts []byte, typ, tit string) (bool, error) {
	defer ruisUtil.Recovers("DataSetParam", func(errs string) {
		anlib.LogError("ruisgo-DataSetParam", "err:"+errs)
	})
	btaes, err := openssl.AesECBEncrypt(bts, core.DbAESKEY, openssl.PKCS7_PADDING)
	//btaes,err:=ruisUtil.AESEncrypt(bts,core.DbAESKEY,core.DbAESIV)
	if err != nil {
		return false, err
	}
	e := DataGetParam(nm)
	if e == nil {
		if len(typ) <= 0 && len(tit) <= 0 {
			anlib.LogError("ruisgo-DataSetParam", "Insert type or title is empty")
			return false, errors.New("type or title is empty")
		}

		e := &DataParam{}
		e.Types = typ
		e.Title = tit
		e.Name = nm
		e.Times = time.Now()
		e.Value = btaes
		_, err := dbMain.Insert(e)
		if err != nil {
			anlib.LogError("ruisgo-DataSetParam", "Insert err:"+err.Error())
			return false, err
		}
	} else {
		clms := make([]string, 0)
		clms = append(clms, "value")
		if len(typ) > 0 && len(tit) > 0 {
			e.Types = typ
			e.Title = tit
			clms = append(clms, "types", "title")
		}

		e.Value = btaes
		_, err := dbMain.Cols(clms...).Where("id=?", e.Id).Update(e)
		if err != nil {
			anlib.LogError("ruisgo-DataSetParam", "Insert err:"+err.Error())
			return false, err
		}
	}

	return true, nil
}
