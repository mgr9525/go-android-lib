package golib

import (
	"errors"
	ruisUtil "github.com/mgr9525/go-ruisutil"
	"github.com/thinkoner/openssl"
	"go-android-lib/anlib"
	"go-android-lib/core"
	"time"
)

type DataUser struct {
	Id     int64  `xorm:"pk autoincr INTEGER NOT NULL"`
	Uid    string `xorm:"unique VARCHAR(50)"`
	Name   string `xorm:"unique VARCHAR(50)"`
	Token  []byte `xorm:"BLOB"`
	Timelg time.Time
}

func (c *DataUser) Tokens() string {
	defer ruisUtil.Recovers("DataUser.Tokens", nil)
	if c.Token == nil || len(c.Token) <= 0 {
		return ""
	}
	//btaes,err:=ruisUtil.AESDecrypt(c.Token,core.DbAESKEY,core.DbAESIV)
	btaes, err := openssl.AesECBDecrypt(c.Token, core.DbAESKEY, openssl.PKCS7_PADDING)
	if err != nil {
		anlib.LogError("ruisgo-DataUser", "Tokens:"+err.Error())
		return ""
	}
	return string(btaes)
}

type DataUserArray struct {
	ls []*DataUser
}

func (c *DataUserArray) Len() int {
	return len(c.ls)
}
func (c *DataUserArray) Get(i int) *DataUser {
	defer ruisUtil.Recovers("DataUserArray.Get", nil)
	return c.ls[i]
}

type DbMainUser struct{}

func (DbMainUser) Find() *DataUserArray {
	defer ruisUtil.Recovers("DbMainUser", func(errs string) {
		anlib.LogError("ruisgo-DbMainUser", "Find:"+errs)
	})
	arr := &DataUserArray{}
	ses := dbMain.NewSession()
	defer ses.Close()
	err := ses.OrderBy("timelg DESC").Find(&arr.ls)
	if err != nil {
		anlib.LogError("ruisgo-FindVDraft", err.Error())
	}
	return arr
}

func (DbMainUser) Get(id int64) *DataUser {
	defer ruisUtil.Recovers("DbMainUser", func(errs string) {
		anlib.LogError("ruisgo-DbMainUser", "Get:"+errs)
	})
	e := &DataUser{}
	ok, err := dbMain.Where("id=?", id).Get(e)
	if err != nil {
		anlib.LogError("ruisgo-GetVDraft", err.Error())
	}
	if ok {
		return e
	}
	return nil
}
func (DbMainUser) LastUser() *DataUser {
	defer ruisUtil.Recovers("DbMainUser", func(errs string) {
		anlib.LogError("ruisgo-DbMainUser", "LastUser:"+errs)
	})
	e := &DataUser{}
	ok, err := dbMain.Where("token is not null").OrderBy("timelg DESC").Get(e)
	if err != nil {
		anlib.LogError("ruisgo-GetVDraft", err.Error())
	}
	if ok {
		return e
	}
	return nil
}
func (DbMainUser) FindUid(xid string) *DataUser {
	defer ruisUtil.Recovers("DbMainUser", func(errs string) {
		anlib.LogError("ruisgo-DbMainUser", "FindUid:"+errs)
	})
	e := &DataUser{}
	ok, err := dbMain.Where("uid=?", xid).Get(e)
	if err != nil {
		anlib.LogError("ruisgo-GetVDraft", err.Error())
	}
	if ok {
		return e
	}
	return nil
}
func (DbMainUser) FindName(name string) *DataUser {
	defer ruisUtil.Recovers("DbMainUser", func(errs string) {
		anlib.LogError("ruisgo-DbMainUser", "FindName:"+errs)
	})
	e := &DataUser{}
	ok, err := dbMain.Where("name=?", name).Get(e)
	if err != nil {
		anlib.LogError("ruisgo-GetVDraft", err.Error())
	}
	if ok {
		return e
	}
	return nil
}
func (DbMainUser) Del(id int64) bool {
	defer ruisUtil.Recovers("DbMainUser", func(errs string) {
		anlib.LogError("ruisgo-DbMainUser", "Del:"+errs)
	})
	if id <= 0 {
		return false
	}
	_, err := dbMain.Where("id=?", id).Delete(DataUser{})
	if err != nil {
		anlib.LogError("ruisgo-DelVDraft", err.Error())
		return false
	}
	return true
}
func (DbMainUser) DelUid(uid string) bool {
	defer ruisUtil.Recovers("DbMainUser", func(errs string) {
		anlib.LogError("ruisgo-DbMainUser", "DelUid:"+errs)
	})
	if len(uid) <= 0 {
		return false
	}
	_, err := dbMain.Where("uid=?", uid).Delete(DataUser{})
	if err != nil {
		anlib.LogError("ruisgo-DelVDraft", err.Error())
		return false
	}
	return true
}
func (c DbMainUser) Login(uid, name, tks string) (int64, error) {
	defer ruisUtil.Recovers("DbMainUser", func(errs string) {
		anlib.LogError("ruisgo-DbMainUser", "Login:"+errs)
	})
	if len(uid) <= 0 || len(name) <= 0 || len(tks) <= 0 {
		return 0, errors.New("param err！")
	}
	isadd := false
	e := c.FindUid(uid)
	if e == nil {
		isadd = true
		e = &DataUser{}
		e.Uid = uid
		e.Name = name
	}
	btaes, erre := openssl.AesECBEncrypt([]byte(tks), core.DbAESKEY, openssl.PKCS7_PADDING)
	//btaes,erre:=ruisUtil.AESEncrypt([]byte(tks),core.DbAESKEY,core.DbAESIV)
	if erre != nil {
		return 0, erre
	}
	e.Timelg = time.Now()
	e.Token = btaes
	var err error
	if isadd {
		_, err = dbMain.Insert(e)
	} else {
		_, err = dbMain.Cols("timelg", "token").Where("id=?", e.Id).Update(e)
	}
	if err != nil {
		anlib.LogError("ruisgo-VDraftAdd", err.Error())
		return 0, err
	}
	return e.Id, nil
}
func (c DbMainUser) Logout() bool {
	defer ruisUtil.Recovers("DbMainUser", func(errs string) {
		anlib.LogError("ruisgo-DbMainUser", "Logout:"+errs)
	})
	usr := c.LastUser()
	if usr != nil {
		usr.Token = nil
		_, err := dbMain.Cols("token").Where("id=?", usr.Id).Update(usr)
		if err != nil {
			anlib.LogError("ruisgo-Logout", err.Error())
			return false
		}
		return true
	}
	return true
}
