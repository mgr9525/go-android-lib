package golib

import (
	ruisUtil "github.com/mgr9525/go-ruisutil"
	"go-android-lib/anlib"
)

type GoRunner interface {
	Run()
}

func NewThread(t GoRunner) {
	go func() {
		defer ruisUtil.Recovers("NewThread", func(errs string) {
			anlib.LogError("ruisgo-NewThread", "run:"+errs)
		})
		t.Run()
	}()
}
