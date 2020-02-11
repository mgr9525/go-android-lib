package golib

import (
	ruisUtil "github.com/mgr9525/go-ruisutil"
	"go-android-lib/anlib"
)

type GoMap struct {
	mp *ruisUtil.Map
}

func NewMap() *GoMap {
	return &GoMap{mp: ruisUtil.NewMap()}
}
func NewMaps(body string) *GoMap {
	return &GoMap{mp: ruisUtil.NewMaps(body)}
}
func NewMapbs(body []byte) *GoMap {
	return &GoMap{mp: ruisUtil.NewMapo(body)}
}
func NewMapo(body *GoMap) *GoMap {
	if body.mp == nil {
		return &GoMap{mp: ruisUtil.NewMap()}
	}
	return &GoMap{mp: ruisUtil.NewMapo(body.mp)}
}
func (c *GoMap) Get(key string) interface{} {
	defer ruisUtil.Recovers("Get", func(errs string) {
		anlib.LogError("ruisgo-GoMap", "Get:"+errs)
	})
	return c.mp.Get(key)
}
func (c *GoMap) Set(key string) interface{} {
	defer ruisUtil.Recovers("Set", func(errs string) {
		anlib.LogError("ruisgo-GoMap", "Set:"+errs)
	})
	return c.mp.Get(key)
}
func (c *GoMap) String() string {
	defer ruisUtil.Recovers("ToString", func(errs string) {
		anlib.LogError("ruisgo-GoMap", "ToString:"+errs)
	})
	return c.mp.ToString()
}
func (c *GoMap) GetString(key string) string {
	defer ruisUtil.Recovers("GetString", func(errs string) {
		anlib.LogError("ruisgo-GoMap", "GetString:"+errs)
	})
	return c.mp.GetString(key)
}
func (c *GoMap) GetInt(key string) (int64, error) {
	defer ruisUtil.Recovers("GetInt", func(errs string) {
		anlib.LogError("ruisgo-GoMap", "GetInt:"+errs)
	})
	return c.mp.GetInt(key)
}
func (c *GoMap) GetFloat(key string) (float64, error) {
	defer ruisUtil.Recovers("GetFloat", func(errs string) {
		anlib.LogError("ruisgo-GoMap", "GetFloat:"+errs)
	})
	return c.mp.GetFloat(key)
}
func (c *GoMap) GetBool(key string) bool {
	defer ruisUtil.Recovers("GetBool", func(errs string) {
		anlib.LogError("ruisgo-GoMap", "GetBool:"+errs)
	})
	return c.mp.GetBool(key)
}
func (c *GoMap) SetString(key, v string) {
	defer ruisUtil.Recovers("SetString", func(errs string) {
		anlib.LogError("ruisgo-GoMap", "SetString:"+errs)
	})
	c.mp.Set(key, v)
}
func (c *GoMap) SetInt(key string, v int64) {
	defer ruisUtil.Recovers("SetString", func(errs string) {
		anlib.LogError("ruisgo-GoMap", "SetString:"+errs)
	})
	c.mp.Set(key, v)
}
func (c *GoMap) SetFloat(key string, v float64) {
	defer ruisUtil.Recovers("SetString", func(errs string) {
		anlib.LogError("ruisgo-GoMap", "SetString:"+errs)
	})
	c.mp.Set(key, v)
}
func (c *GoMap) SetBool(key string, v bool) {
	defer ruisUtil.Recovers("SetString", func(errs string) {
		anlib.LogError("ruisgo-GoMap", "SetString:"+errs)
	})
	c.mp.Set(key, v)
}
func (c *GoMap) SetMap(key string, v *GoMap) {
	defer ruisUtil.Recovers("SetMap", func(errs string) {
		anlib.LogError("ruisgo-GoMap", "SetMap:"+errs)
	})
	c.mp.Set(key, v)
}
func (c *GoMap) SetMaps(key string, v []byte) {
	defer ruisUtil.Recovers("SetMaps", func(errs string) {
		anlib.LogError("ruisgo-GoMap", "SetMaps:"+errs)
	})
	c.mp.Set(key, ruisUtil.NewMapo(v))
}
