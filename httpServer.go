package golib

import (
	"errors"
	"fmt"
	ruisUtil "github.com/mgr9525/go-ruisutil"
	ruisIo "github.com/mgr9525/go-ruisutil/ruisio"
	"go-android-lib/anlib"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var server = &HttpServer{}
var proxyPath = ""

type HttpServer struct {
	port     int
	serv     *http.Server
	CustFuns IHttpServer

	MaxCacheSize int64
}
type IHttpServer interface {
	ProxyNameGenerator(ul string) string
}

func init() {
	server.MaxCacheSize = 1024 * 1024
}
func GetHttpServer() *HttpServer {
	return server
}

var ports = []int{10080, 10081, 10082, 10083, 10084, 10085}

func (c *HttpServer) Start() error {
	if server.serv != nil {
		return errors.New("server already started")
	}
	proxyPath = fmt.Sprintf("%s/cache/proxy", GoApp.SdAppPath())
	anlib.LogInfo("ruisgo-proxyPath", "proxyPath:", proxyPath)
	http.HandleFunc("/proxy/video", c.videoProxyHandle)
	go func() {
		for _, v := range ports {
			c.port = v
			c.serv = &http.Server{Addr: fmt.Sprintf(":%d", v)}
			err := c.serv.ListenAndServe()
			if err == nil {
				break
			} else {
				c.serv = nil
			}
		}
	}()
	return nil
}
func (c *HttpServer) Stop() error {
	if c.serv == nil {
		return errors.New("server not started")
	}
	return c.serv.Shutdown(nil)
}
func (c *HttpServer) Port() int {
	if c.serv == nil {
		return 0
	}
	return c.port
}
func (c *HttpServer) GetProxyUrl(ul string) string {
	return fmt.Sprintf("http://localhost:%d/proxy/video?url=%s", c.port, url.QueryEscape(ul))
}
func (c *HttpServer) videoProxyHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ul := r.FormValue("url")
	name := r.FormValue("name")
	//anlib.LogDebug("ruisgo-VideoProxyHandle","url:"+ul)
	if len(ul) <= 0 {
		w.WriteHeader(404)
		return
	}
	if len(name) > 0 {
		c.videoProxy(name, ul, w, r)
		return
	}
	if c.CustFuns == nil {
		c.videoProxy(ruisUtil.Md5String(ul), ul, w, r)
		return
	}

	c.videoProxy(c.CustFuns.ProxyNameGenerator(ul), ul, w, r)
}

func (c *HttpServer) videoProxy(name, ul string, w http.ResponseWriter, r *http.Request) {
	//path:=GoApp.SdAppPath()
	//anlib.LogDebug("ruisgo-VideoProxyHandle", "name:", name, "，url:"+ul)
	//w.Header().Set("Content-Type","text/html")
	//w.WriteHeader(200)
	nepath := fmt.Sprintf("%s/%s.ok", proxyPath, name)
	if ruisIo.PathExists(nepath) {
		c.distFile(name, nepath, w, r)
	} else {
		c.distUrl(name, ul, w, r)
	}
	go c.rm2CacheSize()
}

var hds1 = []string{"Range"}
var hds2 = []string{"Content-Type", "Cache-Control", "Accept-Ranges", "Content-Range"}
var regRngs = regexp.MustCompile(`^(bytes )(\d*)-(\d*)/(\d+)`)

func (c *HttpServer) distUrl(name, ul string, w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", ul, r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	for _, v := range hds1 {
		hd := r.Header.Get(v)
		if len(hd) > 0 {
			req.Header.Set(v, hd)
		}
	}
	cli := http.DefaultClient
	res, err := cli.Do(req)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer res.Body.Close()

	crngs := res.Header.Get("Content-Range")
	//anlib.LogDebug("ruisgo-distUrl", "name:", name, "，url:"+ul, "，Range:"+r.Header.Get("Range"), "，Content-Range:"+crngs)
	for _, v := range hds2 {
		hd := res.Header.Get(v)
		if len(hd) > 0 {
			w.Header().Set(v, hd)
		}
	}

	start := int64(0)
	//size:=int64(0)
	if len(crngs) > 0 && regRngs.MatchString(crngs) {
		rngs := regRngs.FindAllStringSubmatch(crngs, -1)
		//anlib.LogDebug("ruisgo-distUrl", "rngs1:", rngs[0][1], ",rngs2:", rngs[0][2], ",fmt:", fmt.Sprintf("%v", rngs[0]))
		rg0, err0 := strconv.ParseInt(rngs[0][2], 10, 64)
		if err0 == nil {
			start = rg0
		}
		/*sz,errz:=strconv.ParseInt(rngs[0][4],10,64)
		if errz==nil {
			size=sz
		}*/
	}

	var tmpfl *os.File
	tmpath := fmt.Sprintf("%s/%s.tmp", proxyPath, name)
	nepath := fmt.Sprintf("%s/%s.ok", proxyPath, name)
	if start == 0 {
		os.Remove(tmpath)
		fl, err := os.OpenFile(tmpath, os.O_RDWR|os.O_CREATE, 0644)
		if err == nil {
			defer fl.Close()
			tmpfl = fl
			//anlib.LogDebug("ruisgo-distUrl", "OpenFile ok tmpath:", tmpath)
		} else {
			anlib.LogError("ruisgo-distUrl", "OpenFile tmpfl err:", err.Error())
		}
	}

	w.WriteHeader(res.StatusCode)
	bts := make([]byte, 1024)
	for {
		n, err := res.Body.Read(bts)
		if n > 0 {
			if tmpfl != nil {
				tmpfl.Write(bts[0:n])
			}
			_, errw := w.Write(bts[0:n])
			if errw != nil {
				break
			}
		}
		if err == io.EOF {
			if tmpfl != nil {
				tmpfl.Close()
				os.Rename(tmpath, nepath)
				//anlib.LogInfo("ruisgo-distUrl", "Rename nepath:"+nepath)
			}
			break
		}
	}
}
func (c *HttpServer) distFile(name, nepath string, w http.ResponseWriter, r *http.Request) {
	//anlib.LogInfo("ruisgo-distFile", "OpenFile nepath:"+nepath)
	fl, err := os.OpenFile(nepath, os.O_RDONLY, 0)
	if err != nil {
		anlib.LogInfo("ruisgo-distFile", "OpenFile err:"+err.Error())
		w.WriteHeader(500)
		return
	}
	defer fl.Close()
	w.Header().Add("Content-Type", "video/mpeg4")
	w.Header().Add("Cache-Control", "max-age=360000000")

	stat, err := os.Stat(nepath)
	if err != nil {
		anlib.LogInfo("ruisgo-distFile", "stat file err:"+err.Error())
		w.WriteHeader(500)
		return
	}

	is200 := true
	flsz := stat.Size()
	lnd := flsz
	ranges := r.Header.Get("Range")
	if len(ranges) <= 0 {
		anlib.LogInfo("ruisgo-distFile", "res(200) len:"+strconv.FormatInt(lnd, 10))
	} else {
		is200 = false
		if !strings.HasPrefix(ranges, "bytes=") {
			//anlib.LogInfo("ruisgo-distFile", "range only support bytes!!!")
			w.WriteHeader(500)
			return
		}
		//anlib.LogInfo("ruisgo-distFile", "http ranges:", ranges)
		rgstr := strings.Replace(ranges, "bytes=", "", 1)
		rgs := strings.Split(rgstr, "-")
		if len(rgs) != 2 {
			//anlib.LogInfo("ruisgo-distFile", "range only support bytes err!!!")
			w.WriteHeader(500)
			return
		}
		rgo0 := int64(0)
		rgo1 := int64(0)
		rg0, errg := strconv.ParseInt(rgs[0], 10, 64)
		if errg == nil {
			rgo0 = rg0
			//anlib.LogInfo("ruisgo-distFile", "(206)sl.seek:", strconv.FormatInt(rg0, 10))
			pos, errsk := fl.Seek(rg0, io.SeekStart)
			if errsk != nil {
				anlib.LogInfo("ruisgo-distFile", "(206)sl.seek err:", errsk.Error())
			}
			anlib.LogInfo("ruisgo-distFile", "(206)sl.seek pos:", strconv.FormatInt(pos, 10))
		}
		rg1, errg := strconv.ParseInt(rgs[1], 10, 64)
		if errg == nil {
			rgo1 = rg1
		}
		if rgo1 <= 0 {
			rgo1 = flsz - 1
		}
		lnd = rgo1 - rgo0 + 1
		w.Header().Add("Accept-Ranges", "bytes")
		w.Header().Add("Content-Range", fmt.Sprintf("bytes %d-%d/%d", rgo0, rgo1, flsz))
		//c.Resp.Header().Add("Content-Length", fmt.Sprintf("%d", lnd))
	}

	w.Header().Add("Content-Length", fmt.Sprintf("%d", lnd))
	if is200 {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(206)
	}
	//anlib.LogInfo("ruisgo-distFile", "will send len:", strconv.FormatInt(lnd, 10))
	if lnd <= 0 {
		return
	}

	bts := make([]byte, 1024)
	for {
		n, err := fl.Read(bts)
		if n > 0 {
			if lnd < int64(n) {
				w.Write(bts[:lnd])
				break
			}
			_, errw := w.Write(bts[:n])
			if errw != nil {
				anlib.LogError("ruisgo-distFile", "resp write err:", errw.Error())
				break
			}
			lnd -= int64(n)
		}
		if err == io.EOF {
			break
		}
	}
	//anlib.LogInfo("ruisgo-distFile", "other send len:", fmt.Sprintf("%d", lnd))
}
func (c *HttpServer) rm2CacheSize() {
	defer ruisUtil.Recovers("HttpServer", func(errs string) {
		anlib.LogError("ruisgo-HttpServer", "rm2CacheSize err:", errs)
	})

	for c.serv != nil {
		zsz := int64(0)
		var lastfl *os.FileInfo
		filepath.Walk(proxyPath, func(path string, info os.FileInfo, err error) error {
			zsz += info.Size()
			if lastfl == nil {
				lastfl = &info
			} else if info.ModTime().Sub((*lastfl).ModTime()).Microseconds() < 0 {
				lastfl = &info
			}
			return err
		})

		if zsz < c.MaxCacheSize {
			break
		}

		rmpath := fmt.Sprintf("%s/%s", proxyPath, (*lastfl).Name())
		anlib.LogDebug("ruisgo-rm2CacheSize", "proxyPath stat:", fmt.Sprintf("size:%d,max:%d,rm file:%s", zsz, c.MaxCacheSize, rmpath))
		err := os.Remove(rmpath)
		if err != nil {
			anlib.LogError("ruisgo-rm2CacheSize", fmt.Sprintf("rm file(%s) err:%s", rmpath, err.Error()))
			break
		}
	}
}
