package httpx

import (
	"github.com/qinchende/gofast/cst"
)

const XForwardFor = "X-Forward-For"

const (
	emptyJson         = "{}"
	maxMemory         = 32 << 20 // 32MB
	maxBodyLen        = 8 << 20  // 8MB
	separator         = ";"
	tokensInAttribute = 2
)

const (
	FormatJson = iota
	FormatUrlEncoding
	FormatXml
)

type RequestPet struct {
	Method     string
	Url        string
	Headers    cst.WebKV
	QueryArgs  cst.KV
	BodyArgs   cst.KV
	BodyFormat int8
}