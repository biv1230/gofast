package logx

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const gftSupportMinGoVer = 10

var isDebug = false
var _once sync.Once

// 每个程序只设置一次debug标志位，后面的设置都失效
func SetDebugStatus(yn bool) {
	_once.Do(func() {
		isDebug = yn
	})
}

func IsDebugging() bool {
	return isDebug
}

// 可以自定义 logx.DebugPrintRouteFunc 改变Debug输出的路由信息
var DebugPrintRouteFunc func(httpMethod, absolutePath, handlerName string, nuHandlers int)

//func DebugPrintRoute(httpMethod, absolutePath string, hds fst.CtxHandlers) {
//	if IsDebugging() {
//		nuHandlers := len(hds)
//		handlerName := lang.NameOfFunc(hds.Last())
//		if DebugPrintRouteFunc == nil {
//			DebugPrint("%-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
//		} else {
//			DebugPrintRouteFunc(httpMethod, absolutePath, handlerName, nuHandlers)
//		}
//	}
//}

func DebugPrint(format string, values ...interface{}) {
	if isDebug {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(os.Stdout, "[GoFast-debug] "+format, values...)
	}
}

func DebugPrintError(err error) {
	if err != nil && isDebug {
		fmt.Fprintf(os.Stderr, "[GoFast-debug] [ERROR] %v\n", err)
	}
}

func GetMinVer(v string) (uint64, error) {
	first := strings.IndexByte(v, '.')
	last := strings.LastIndexByte(v, '.')
	if first == last {
		return strconv.ParseUint(v[first+1:], 10, 64)
	}
	return strconv.ParseUint(v[first+1:last], 10, 64)
}

func DebugPrintWARNINGDefault() {
	if v, e := GetMinVer(runtime.Version()); e == nil && v <= gftSupportMinGoVer {
		DebugPrint(`[WARNING] Now GoFast requires Go 1.11 or later and Go 1.12 will be required soon.

`)
	}
	DebugPrint(`[WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

`)
}

//func DebugPrintWARNINGNew() {
//	DebugPrint(`[WARNING] Running in "debug" mode. Switch to "release" mode in production.
// - using env:	export GoFast_MODE=release
// - using code:	fst.SetMode(fst.ReleaseMode)
//
//`)
//}
//
//func DebugPrintWARNINGSetHTMLTemplate() {
//	DebugPrint(`[WARNING] Since SetHTMLTemplate() is NOT thread-safe. It should only be called
//at initialization. ie. before any route is registered or the router is listening in a socket:
//
//	router := fst.Default()
//	router.SetHTMLTemplate(template) // << good place
//
//`)
//}