package golib

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	ruisUtil "github.com/mgr9525/go-ruisutil"
	"go-android-lib/anlib"
	"go-android-lib/core"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpCallback interface {
	Ok(res *Response)
	Err(errs error)
}
type Response struct {
	res *http.Response
}

func (c *Response) StatusCode() int {
	defer ruisUtil.Recovers("StatusCode", func(errs string) {
		anlib.LogError("ruisgo-Response", "StatusCode:"+errs)
	})
	return c.res.StatusCode
}
func (c *Response) GetHeader(k string) string {
	defer ruisUtil.Recovers("GetHeader", func(errs string) {
		anlib.LogError("ruisgo-Response", "GetHeader:"+errs)
	})
	return c.res.Header.Get(k)
}
func (c *Response) ContentLen() int64 {
	defer ruisUtil.Recovers("ContentLen", func(errs string) {
		anlib.LogError("ruisgo-Response", "ContentLen:"+errs)
	})
	return c.res.ContentLength
}
func (c *Response) BodyRead(bts []byte) (int, error) {
	defer ruisUtil.Recovers("BodyRead", func(errs string) {
		anlib.LogError("ruisgo-Response", "BodyRead:"+errs)
	})

	return c.res.Body.Read(bts)
}
func (c *Response) BodyClose(bts []byte) error {
	defer ruisUtil.Recovers("BodyClose", func(errs string) {
		anlib.LogError("ruisgo-Response", "BodyClose:"+errs)
	})
	return c.res.Body.Close()
}
func (c *Response) BodyReadAll() ([]byte, error) {
	defer ruisUtil.Recovers("BodyReadAll", func(errs string) {
		anlib.LogError("ruisgo-Response", "BodyReadAll:"+errs)
	})
	bts, err := ioutil.ReadAll(c.res.Body)
	if err != nil {
		return nil, err
	}
	return bts, nil
}
func (c *Response) BodyReadAlls() (string, error) {
	defer ruisUtil.Recovers("BodyReadAlls", func(errs string) {
		anlib.LogError("ruisgo-Response", "BodyReadAlls:"+errs)
	})
	bts, err := c.BodyReadAll()
	if err != nil {
		return "", err
	}
	return string(bts), nil
}
func (c *Response) BodyReadAllo() (*GoMap, error) {
	defer ruisUtil.Recovers("BodyReadAllo", func(errs string) {
		anlib.LogError("ruisgo-Response", "BodyReadAllo:"+errs)
	})
	bts, err := c.BodyReadAll()
	if err != nil {
		return nil, err
	}
	return NewMapbs(bts), nil
}

type Request struct {
	req     *http.Request
	Method  string
	Url     string
	Header  *GoMap
	Timeout int64 //ms

	Isapp bool
	Isusr bool
}

func NewHttpReq(ul string) *Request {
	return &Request{Url: ul, Header: NewMap()}
}
func NewHttpApi(ul string) *Request {
	return &Request{Url: ul, Header: NewMap(), Isusr: true, Isapp: true}
}
func getAppToken() string {
	cont := fmt.Sprintf("%s|%s", goAppXid, time.Now().Format(time.RFC3339Nano))
	sign := ruisUtil.Sha1String(fmt.Sprintf("%s.%s", cont, core.AppKey))
	s := base64.StdEncoding.EncodeToString([]byte(cont))
	return fmt.Sprintf("%s.%s", s, sign)
}
func (c *Request) doGet() (ret *Response, reterr error) {
	defer ruisUtil.Recovers("doGet", func(errs string) {
		reterr = errors.New(errs)
		anlib.LogError("ruisgo-Request", "doGet:"+errs)
	})
	req, err := http.NewRequest(c.Method, c.Url, nil)
	if err != nil {
		return nil, err
	}
	if c.Header != nil && c.Header.mp != nil {
		for k, _ := range c.Header.mp.Map() {
			req.Header.Set(k, c.Header.GetString(k))
		}
	}
	if c.Isapp {
		req.Header.Set("app-token", getAppToken())
	}
	if c.Isusr {
		usr := DbMainUser{}.LastUser()
		if usr != nil {
			req.Header.Set("Authorization", usr.Tokens())
		}
	}
	cli := &http.Client{}
	if c.Timeout > 0 {
		cli.Timeout = time.Duration(c.Timeout) * time.Millisecond
	}
	res, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	rets := &Response{}
	rets.res = res
	return rets, nil
}
func (c *Request) doPost(body *GoMap) (ret *Response, reterr error) {
	defer ruisUtil.Recovers("doPost", func(errs string) {
		reterr = errors.New(errs)
		anlib.LogError("ruisgo-Request", "doPost:"+errs)
	})
	pars := &url.Values{}
	if body != nil && body.mp != nil {
		for k, _ := range body.mp.Map() {
			pars.Set(k, body.GetString(k))
		}
	}
	req, err := http.NewRequest(c.Method, c.Url, strings.NewReader(pars.Encode()))
	if err != nil {
		return nil, err
	}
	if c.Header != nil && c.Header.mp != nil {
		for k, _ := range c.Header.mp.Map() {
			req.Header.Set(k, c.Header.GetString(k))
		}
	}
	if c.Isapp {
		req.Header.Set("app-token", getAppToken())
	}
	if c.Isusr {
		usr := DbMainUser{}.LastUser()
		if usr != nil {
			req.Header.Set("Authorization", usr.Tokens())
		}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	cli := &http.Client{}
	if c.Timeout > 0 {
		cli.Timeout = time.Duration(c.Timeout) * time.Millisecond
	}
	res, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	rets := &Response{}
	rets.res = res
	return rets, nil
}
func (c *Request) doPostJson(body []byte) (ret *Response, reterr error) {
	defer ruisUtil.Recovers("doPostJson", func(errs string) {
		reterr = errors.New(errs)
		anlib.LogError("ruisgo-Request", "doPostJson:"+errs)
	})
	req, err := http.NewRequest(c.Method, c.Url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	if c.Header != nil && c.Header.mp != nil {
		for k, _ := range c.Header.mp.Map() {
			req.Header.Set(k, c.Header.GetString(k))
		}
	}
	if c.Isapp {
		req.Header.Set("app-token", getAppToken())
	}
	if c.Isusr {
		usr := DbMainUser{}.LastUser()
		if usr != nil {
			req.Header.Set("Authorization", usr.Tokens())
		}
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	cli := &http.Client{}
	if c.Timeout > 0 {
		cli.Timeout = time.Duration(c.Timeout) * time.Millisecond
	}
	res, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	rets := &Response{}
	rets.res = res
	return rets, nil
}
func (c *Request) Do(body *GoMap) (ret *Response, reterr error) {
	defer ruisUtil.Recovers("Do", func(errs string) {
		anlib.LogError("ruisgo-Request", "Do:"+errs)
		reterr = errors.New(errs)
	})
	if len(c.Method) <= 0 {
		if body == nil {
			c.Method = "GET"
		} else {
			c.Method = "POST"
		}
	}
	mth := strings.ToUpper(c.Method)
	if mth == "GET" {
		return c.doGet()
	} else if mth == "POST" {
		return c.doPost(body)
	}
	return nil, errors.New("not support method")
}
func (c *Request) DoJson(body string) (ret *Response, reterr error) {
	defer ruisUtil.Recovers("DoJson", func(errs string) {
		anlib.LogError("ruisgo-Request", "DoJsons:"+errs)
		reterr = errors.New(errs)
	})
	c.Method = "POST"
	return c.doPostJson([]byte(body))
}
func (c *Request) DoJsons(body []byte) (ret *Response, reterr error) {
	defer ruisUtil.Recovers("DoJson", func(errs string) {
		anlib.LogError("ruisgo-Request", "DoJsons:"+errs)
		reterr = errors.New(errs)
	})
	c.Method = "POST"
	return c.doPostJson(body)
}
func (c *Request) DoJsono(body *GoMap) (ret *Response, reterr error) {
	defer ruisUtil.Recovers("DoJsono", func(errs string) {
		anlib.LogError("ruisgo-Request", "DoJsono:"+errs)
		reterr = errors.New(errs)
	})
	c.Method = "POST"
	bts, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return c.doPostJson(bts)
}
func (c *Request) DoAsyn(body *GoMap, callback HttpCallback) {
	go func() {
		defer ruisUtil.Recovers("DoAsyn", func(errs string) {
			anlib.LogError("ruisgo-Request", "DoAsyn:"+errs)
			callback.Err(errors.New(errs))
		})

		res, err := c.Do(body)
		if err == nil {
			callback.Ok(res)
		} else {
			callback.Err(err)
		}
	}()
}
func (c *Request) DoJsonAsyn(body string, callback HttpCallback) {
	bts := []byte(body)
	bodys := make([]byte, len(bts))
	copy(bodys, bts)
	go func() {
		defer ruisUtil.Recovers("DoJsonAsyn", func(errs string) {
			anlib.LogError("ruisgo-Request", "DoJsonAsyn:"+errs)
			callback.Err(errors.New(errs))
		})

		res, err := c.DoJson(string(bodys))
		if err == nil {
			callback.Ok(res)
		} else {
			callback.Err(err)
		}
	}()
}
func (c *Request) DoJsonsAsyn(body []byte, callback HttpCallback) {
	bodys := make([]byte, len(body))
	copy(bodys, body)
	go func() {
		defer ruisUtil.Recovers("DoJsonsAsyn", func(errs string) {
			anlib.LogError("ruisgo-Request", "DoJsonsAsyn:"+errs)
			callback.Err(errors.New(errs))
		})

		res, err := c.DoJsons(bodys)
		if err == nil {
			callback.Ok(res)
		} else {
			callback.Err(err)
		}
	}()
}
func (c *Request) DoJsonoAsyn(body *GoMap, callback HttpCallback) {
	go func() {
		defer ruisUtil.Recovers("DoJsonoAsyn", func(errs string) {
			anlib.LogError("ruisgo-Request", "DoJsonoAsyn:"+errs)
			callback.Err(errors.New(errs))
		})

		res, err := c.DoJsono(body)
		if err == nil {
			callback.Ok(res)
		} else {
			callback.Err(err)
		}
	}()
}
