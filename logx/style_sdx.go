// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package logx

import (
	"fmt"
	"github.com/qinchende/gofast/fst/tools"
	"github.com/qinchende/gofast/skill/jsonx"
	"github.com/qinchende/gofast/skill/lang"
	"github.com/qinchende/gofast/skill/timex"
	"strconv"
	"strings"
	"time"
)

func outputSdxStyle(w WriterCloser, logLevel string, data any) {
	info := lang.ToString(data)
	sb := strings.Builder{}
	sb.Grow(len(info) + 26)
	sb.WriteByte('[')
	sb.WriteString(time.Now().Format(timeFormatMini))
	sb.WriteString("][")
	sb.WriteString(logLevel)
	sb.WriteString("]: ")
	sb.WriteString(info)
	outputDirectBuilder(w, &sb)
}

// 通过模板构造字符串可能性能更好。
func buildSdxReqLog(p *ReqLogEntity, flag int8) string {
	// 需要用Mini版本
	if flag > 0 {
		return buildSdxReqLogMini(p)
	}

	formatStr := `
[%s] %s (%s/%s) [%d/%d/%d]
  B: %s
  P: %s
  R: %s%s
`
	// 最长打印出 1024个字节的结果
	tLen := p.BodySize
	if tLen > 1024 {
		tLen = 1024
	}

	// 这个时候可以随意改变 p.Pms ，这是请求最后一个执行的地方了
	var basePms = make(map[string]any)
	if p.Pms["tok"] != nil {
		basePms["tok"] = p.Pms["tok"]
		delete(p.Pms, "tok")
	}

	// 请求参数
	var reqParams []byte
	if p.Pms != nil {
		reqParams, _ = jsonx.Marshal(p.Pms)
	} else if p.RawReq.Form != nil {
		reqParams, _ = jsonx.Marshal(p.RawReq.Form)
	}
	// 请求 核心参数
	reqBaseParams, _ := jsonx.Marshal(basePms)

	return fmt.Sprintf(formatStr,
		p.RawReq.Method,
		p.RawReq.URL.Path,
		p.ClientIP,
		timex.ToTime(p.TimeStamp).Format(timeFormatMini),
		p.StatusCode,
		p.BodySize,
		p.Latency/time.Millisecond,
		reqBaseParams,
		reqParams,
		(p.ResData)[:tLen],
		logBaskets(p.MsgBaskets),
	)
}

func buildSdxReqLogMini(p *ReqLogEntity) string {
	formatStr := `
[%s] %s (%s/%s) [%d/%d/%d] %s
`
	// 最长打印出 1024个字节的结果
	tLen := p.BodySize
	if tLen > 1024 {
		tLen = 1024
	}

	return fmt.Sprintf(formatStr,
		p.RawReq.Method,
		p.RawReq.URL.Path,
		p.ClientIP,
		timex.ToTime(p.TimeStamp).Format(timeFormatMini),
		p.StatusCode,
		p.BodySize,
		p.Latency/time.Millisecond,
		(p.ResData)[:tLen],
	)
}

// 所有错误合并成字符串
func logBaskets(bs tools.Baskets) string {
	if len(bs) == 0 {
		return ""
	}

	var buf strings.Builder
	buf.Grow(len(bs[0].Msg) + 10)

	buf.WriteString("\n  E: ")
	infos := bs.CollectMessages()
	for i, str := range infos {
		if i != 0 {
			buf.WriteString("\n     ")
		}
		buf.WriteString(strconv.Itoa(i))
		buf.WriteString(". ")
		buf.WriteString(str)
	}
	return buf.String()
}

//func traceInfo(ctx context.Context) string {
//	jsonx.m
//
//	return jsonx.UnmarshalFromString() trace.SpanContextFromContext(ctx)
//}
