package golib

import (
	"crypto/rand"
	"errors"
	"fmt"
	ruisUtil "github.com/mgr9525/go-ruisutil"
	ruisIo "github.com/mgr9525/go-ruisutil/ruisio"
	"go-android-lib/anlib"
	"io"
	"net/http"
	"os"
	"time"
)

type UploadListen interface {
	Progress(send, size int64, str string)
	Err(err error)
	Ok(res *Response)
}

type Uploader struct {
	isRun  bool
	Url    string
	Path   string
	Header *GoMap

	listen UploadListen

	Isapp bool
	Isusr bool

	client        *http.Client
	req           *http.Request
	param         *GoMap
	body          *ruisUtil.CircleByteBuffer
	boundarybytes []byte
	endbytes      []byte
	canends       bool
}

func NewUploader(ul, fl string, ltn UploadListen) *Uploader {
	if len(ul) <= 0 || len(fl) <= 0 {
		return nil
	}
	return &Uploader{Url: ul, Path: fl, listen: ltn}
}
func (c *Uploader) Start(pars *GoMap) error {
	if c.isRun {
		return errors.New("aleary started")
	}
	if !ruisIo.PathExists(c.Path) {
		anlib.LogError("ruisgo-Uploader", "Start:file(path:", c.Path, ") not exist")
		return errors.New("file(path) not exist")
	}
	anlib.LogDebug("ruisgo-Uploader", "start init")
	err := c.init()
	if err != nil {
		return err
	}
	c.param = pars
	c.isRun = true
	go c.run()
	return nil
}
func (c *Uploader) Stop() {
	defer ruisUtil.Recovers("Uploader.stop", nil)
	c.isRun = false
	c.body.Close()
	if c.client != nil {
		c.client.CloseIdleConnections()
	}
}
func (c *Uploader) IsRun() bool {
	return c.isRun
}

func (c *Uploader) randomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}
func (c *Uploader) init() error {
	boundary := c.randomBoundary()
	c.body = ruisUtil.NewCircleByteBuffer(1024 * 10)
	c.boundarybytes = []byte("\r\n--" + boundary + "\r\n")
	c.endbytes = []byte("\r\n--" + boundary + "--\r\n")

	reqest, err := http.NewRequest("POST", c.Url, c.body)
	if err != nil {
		return err
	}
	reqest.Header.Add("Connection", "keep-alive")
	reqest.Header.Add("Content-Type", "multipart/form-data; charset=utf-8; boundary="+boundary)
	if c.Isapp {
		reqest.Header.Set("app-token", getAppToken())
	}
	if c.Isusr {
		usr := DbMainUser{}.LastUser()
		if usr != nil {
			reqest.Header.Set("Authorization", usr.Tokens())
		}
	}
	c.req = reqest
	return nil
}
func (c *Uploader) run() {
	defer ruisUtil.Recovers("Uploader.run", func(errs string) {
		anlib.LogError("ruisgo-Uploader", "run:"+errs)
	})

	anlib.LogDebug("ruisgo-Uploader", "http Do Start!")
	go c.send()
	c.client = http.DefaultClient
	res, err := c.client.Do(c.req)
	anlib.LogDebug("ruisgo-Uploader", "http Do End!")
	if err != nil {
		c.listen.Err(err)
		return
	}
	c.listen.Ok(&Response{res: res})
	c.isRun = false
	c.body.Close()
}

func (c *Uploader) send() {
	defer ruisUtil.Recovers("Uploader.send", func(errs string) {
		anlib.LogError("ruisgo-Uploader", "send:"+errs)
	})
	anlib.LogDebug("ruisgo-Uploader", "send() start")
	c.canends = false
	defer func() {
		if c.canends {
			c.body.Write(c.endbytes)
		}
		c.body.Write(nil)
		anlib.LogDebug("ruisgo-Uploader", "send defer!")
	}()
	f, err := os.OpenFile(c.Path, os.O_RDONLY, 0666)
	if err != nil {
		c.listen.Err(err)
		return
	}
	stat, err := f.Stat()
	if err != nil {
		c.listen.Err(err)
		return
	}
	defer f.Close()

	if c.param != nil {
		for k, _ := range c.param.mp.Map() {
			des := fmt.Sprintf("Content-Disposition: form-data; name=\"%s\"\r\n\r\n", k)
			c.body.Write(c.boundarybytes)
			c.body.Write([]byte(des))
			c.body.Write([]byte(c.param.GetString(k)))
		}
	}

	anlib.LogDebug("ruisgo-Uploader", "send() write")
	des := fmt.Sprintf("Content-Disposition: form-data; name=\"upfile\"; filename=\"%s\"\r\nContent-Type: application/octet-stream\r\n\r\n", stat.Name())
	c.body.Write(c.boundarybytes)
	c.body.Write([]byte(des))

	//c.body.Write([]byte("hahsdfhasdfhsahdlkfjsalkjewoirjoweiurowijflskajfdl!!!!!!!!!!!!"))

	fsz := stat.Size()
	fupsz := int64(0)
	buf := make([]byte, 1024)
	for c.isRun {
		time.Sleep(10 * time.Microsecond)
		n, err := f.Read(buf)
		if n > 0 {
			nz, _ := c.body.Write(buf[0:n])
			fupsz += int64(nz)
			xs := float64(fupsz) / float64(fsz)
			c.listen.Progress(fupsz, fsz, fmt.Sprintf("%.2f", xs*100))
		}
		if err == io.EOF {
			c.canends = true
			break
		}
	}
}
