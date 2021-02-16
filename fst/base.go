// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"fmt"
	"github.com/qinchende/gofast/fst/binding"
	"io"
	"math"
	"net/http"
	"os"
)

type (
	I           interface{}
	KV          map[string]interface{}
	IncHandler  func(w http.ResponseWriter, r *Request)
	IncHandlers []IncHandler
	CtxHandler  func(ctx *Context)
	CtxHandlers []CtxHandler
	AppHandler  func(gft *GoFast)
	AppHandlers []AppHandler
)

const (
	// BodyBytesKey indicates a default body bytes key.
	BodyBytesKey     = "_qinchende/gofast/bodybyteskey"
	maxFitLen    int = math.MaxInt8 // 最多多少个中间件函数
	//routePathMaxLen    uint8 = 255      // 路由字符串最长长度
	//routeMaxHandlers   uint8 = 255      // 路由 handlers 最大长度
	defMultipartMemory int64 = 32 << 20 // 32 MB
)

// Content-Type MIME of the most common data formats.
// 常量值，供外部访问调用
const (
	MIMEJSON              = binding.MIMEJSON
	MIMEHTML              = binding.MIMEHTML
	MIMEXML               = binding.MIMEXML
	MIMEXML2              = binding.MIMEXML2
	MIMEPlain             = binding.MIMEPlain
	MIMEPOSTForm          = binding.MIMEPOSTForm
	MIMEMultipartPOSTForm = binding.MIMEMultipartPOSTForm
	MIMEYAML              = binding.MIMEYAML
)

var (
	spf            = fmt.Sprintf
	mimePlain      = []string{MIMEPlain}
	default404Body = []byte("404 (PAGE NOT FOND)")
	default405Body = []byte("405 (METHOD NOT ALLOWED)")

	DefaultWriter      io.Writer = os.Stdout
	DefaultErrorWriter io.Writer = os.Stderr
)
