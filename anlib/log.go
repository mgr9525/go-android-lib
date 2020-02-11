package anlib

/*
#cgo LDFLAGS: -landroid -llog

#include <stdlib.h>
#include <android/log.h>

void anlib_logs(int fs,const char*tag,const char*conts);
*/
import "C"
import (
	"go-android-lib/core"
	"strings"
	"unsafe"
)

func Log(fs int, tag, cont string) {
	tags := C.CString(tag)
	conts := C.CString(cont)
	C.anlib_logs(C.int(fs), tags, conts)
	C.free(unsafe.Pointer(tags))
	C.free(unsafe.Pointer(conts))
}

func LogInfo(tag string, conts ...string) {
	Log(core.LOG_INFO, tag, strings.Join(conts, "  "))
}
func LogDebug(tag string, conts ...string) {
	Log(core.LOG_DEBUG, tag, strings.Join(conts, "  "))
}
func LogError(tag string, conts ...string) {
	Log(core.LOG_ERROR, tag, strings.Join(conts, "  "))
}
