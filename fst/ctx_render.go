// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"errors"
	"fmt"
	"github.com/qinchende/gofast/fst/binding"
	"github.com/qinchende/gofast/fst/render"
	"github.com/qinchende/gofast/fst/sse"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// GoFast JSON render
// JSON是GoFast默认的返回格式，一等公民

func (c *Context) Fai(msg string) {
	c.FaiX(0, msg, nil)
}

func (c *Context) FaiErr(err error) {
	c.FaiX(0, err.Error(), nil)
}

func (c *Context) FaiX(code int32, msg string, obj interface{}) {
	jsonData := KV{
		"status":   "fai",
		"msg_code": code,
		"msg":      msg,
	}
	if obj != nil {
		jsonData["record"] = obj
	}
	c.FaiKV(jsonData)
}

func (c *Context) FaiKV(jsonData KV) {
	if jsonData == nil {
		jsonData = make(KV)
	}
	jsonData["status"] = "fai"
	if jsonData["msg"] == nil {
		jsonData["msg"] = ""
	}
	if jsonData["msg_code"] == nil {
		jsonData["msg_code"] = 0
	}
	if c.Sess.IsNew {
		jsonData["tok"] = c.Sess.Token
	}

	c.JSONOfKV(http.StatusOK, jsonData)
	c.aborted = true
}

// ++
func (c *Context) Suc(obj interface{}) {
	c.SucX(0, "success", obj)
}

func (c *Context) SucMsg(msg string, obj interface{}) {
	c.SucX(0, msg, obj)
}

func (c *Context) SucX(code int32, msg string, obj interface{}) {
	jsonData := KV{
		"status":   "suc",
		"msg_code": code,
		"msg":      msg,
	}
	if obj != nil {
		jsonData["record"] = obj
	}
	c.SucKV(jsonData)
}

func (c *Context) SucKV(jsonData KV) {
	if jsonData == nil {
		jsonData = make(KV)
	}
	jsonData["status"] = "suc"
	if jsonData["msg"] == nil {
		jsonData["msg"] = ""
	}
	if jsonData["msg_code"] == nil {
		jsonData["msg_code"] = 0
	}
	if c.Sess.IsNew {
		jsonData["tok"] = c.Sess.Token
	}

	c.JSONOfKV(http.StatusOK, jsonData)
	c.aborted = true
}

/************************************/
/******** RESPONSE RENDERING ********/
/************************************/

// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}

// Status sets the HTTP response code.
func (c *Context) Status(code int) {
	c.ResW.WriteHeader(code)
}

// Header is a intelligent shortcut for c.ResW.Header().Set(key, value).
// It writes a header in the response.
// If value == "", this method removes the header `c.ResW.Header().Del(key)`
func (c *Context) Header(key, value string) {
	if value == "" {
		c.ResW.Header().Del(key)
		return
	}
	c.ResW.Header().Set(key, value)
}

// GetHeader returns value from request headers.
func (c *Context) GetHeader(key string) string {
	return c.requestHeader(key)
}

// GetRawData return stream data.
func (c *Context) GetRawData() ([]byte, error) {
	return ioutil.ReadAll(c.ReqW.Body)
}

// SetSameSite with cookie
func (c *Context) SetSameSite(samesite http.SameSite) {
	c.sameSite = samesite
}

// SetCookie adds a Set-Cookie header to the ResponseWriter's headers.
// The provided cookie must have a valid Name. Invalid cookies may be
// silently dropped.
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.ResW, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: c.sameSite,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

// Cookie returns the named cookie provided in the request or
// ErrNoCookie if not found. And return the named cookie is unescaped.
// If multiple cookies match the given name, only one cookie will
// be returned.
func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.ReqW.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Render writes the response headers and calls render.Render to render data.
// 返回数据的接口
// 如果需要
func (c *Context) Render(code int, r render.Render) {
	// Render之前加入对应的 render 数据
	c.PRender = &r
	c.PCode = &code

	// add preSend & afterSend events by sdx on 2021.01.06
	c.execPreSendHandlers()

	c.Status(code)
	if !bodyAllowedForStatus(code) {
		r.WriteContentType(c.ResW)
		c.ResW.WriteHeaderNow()
		return
	}
	// TODO: 返回结果之前，统一保存 session
	if c.Sess != nil {
		c.Sess.Save()
	}

	// TODO: 是否要避免 double render，这里的Render是否只需要调一次，如果第二次就需要报错
	if err := r.Render(c.ResW); err != nil {
		panic(err)
	}
	// add preSend & afterSend events by sdx on 2021.01.06
	c.execAfterSendHandlers()
}

// HTML renders the HTTP template specified by its file name.
// It also updates the HTTP code and sets the Content-Type as "text/html".
// See http://golang.org/doc/articles/wiki/
func (c *Context) HTML(code int, name string, obj interface{}) {
	instance := c.gftApp.HTMLRender.Instance(name, obj)
	c.Render(code, instance)
}

// IndentedJSON serializes the given struct as pretty JSON (indented + endlines) into the response body.
// It also sets the Content-Type as "application/json".
// WARNING: we recommend to use this only for development purposes since printing pretty JSON is
// more CPU and bandwidth consuming. Use Context.JSON() instead.
func (c *Context) IndentedJSON(code int, obj interface{}) {
	c.Render(code, render.IndentedJSON{Data: obj})
}

// SecureJSON serializes the given struct as Secure JSON into the response body.
// Default prepends "while(1)," to response body if the given struct is array values.
// It also sets the Content-Type as "application/json".
func (c *Context) SecureJSON(code int, obj interface{}) {
	c.Render(code, render.SecureJSON{Prefix: c.gftApp.SecureJsonPrefix, Data: obj})
}

// JSONP serializes the given struct as JSON into the response body.
// It add padding to response body to request data from a server residing in a different domain than the client.
// It also sets the Content-Type as "application/javascript".
func (c *Context) JSONP(code int, obj interface{}) {
	callback := c.DefaultQuery("callback", "")
	if callback == "" {
		c.Render(code, render.JSON{Data: obj})
		return
	}
	c.Render(code, render.JsonpJSON{Callback: callback, Data: obj})
}

func (c *Context) JSONOfKV(code int, obj KV) {
	c.Render(code, render.JSON{Data: obj})
}

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (c *Context) JSON(code int, obj interface{}) {
	c.Render(code, render.JSON{Data: obj})
}

// AsciiJSON serializes the given struct as JSON into the response body with unicode to ASCII string.
// It also sets the Content-Type as "application/json".
func (c *Context) AsciiJSON(code int, obj interface{}) {
	c.Render(code, render.AsciiJSON{Data: obj})
}

// PureJSON serializes the given struct as JSON into the response body.
// PureJSON, unlike JSON, does not replace special html characters with their unicode entities.
func (c *Context) PureJSON(code int, obj interface{}) {
	c.Render(code, render.PureJSON{Data: obj})
}

// XML serializes the given struct as XML into the response body.
// It also sets the Content-Type as "application/xml".
func (c *Context) XML(code int, obj interface{}) {
	c.Render(code, render.XML{Data: obj})
}

// YAML serializes the given struct as YAML into the response body.
func (c *Context) YAML(code int, obj interface{}) {
	c.Render(code, render.YAML{Data: obj})
}

// ProtoBuf serializes the given struct as ProtoBuf into the response body.
func (c *Context) ProtoBuf(code int, obj interface{}) {
	c.Render(code, render.ProtoBuf{Data: obj})
}

// String writes the given string into the response body.
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Render(code, render.String{Format: format, Data: values})
}

// Redirect returns a HTTP redirect to the specific location.
func (c *Context) Redirect(code int, location string) {
	c.Render(-1, render.Redirect{
		Code:     code,
		Location: location,
		Request:  c.ReqW,
	})
}

// Data writes some data into the body stream and updates the HTTP code.
func (c *Context) Data(code int, contentType string, data []byte) {
	c.Render(code, render.Data{
		ContentType: contentType,
		Data:        data,
	})
}

// DataFromReader writes the specified reader into the body stream and updates the HTTP code.
func (c *Context) DataFromReader(code int, contentLength int64, contentType string, reader io.Reader, extraHeaders map[string]string) {
	c.Render(code, render.Reader{
		Headers:       extraHeaders,
		ContentType:   contentType,
		ContentLength: contentLength,
		Reader:        reader,
	})
}

// File writes the specified file into the body stream in a efficient way.
func (c *Context) File(filepath string) {
	http.ServeFile(c.ResW, c.ReqW, filepath)
}

// FileFromFS writes the specified file from http.FileSytem into the body stream in an efficient way.
func (c *Context) FileFromFS(filepath string, fs http.FileSystem) {
	defer func(old string) {
		c.ReqW.URL.Path = old
	}(c.ReqW.URL.Path)

	c.ReqW.URL.Path = filepath

	http.FileServer(fs).ServeHTTP(c.ResW, c.ReqW)
}

// FileAttachment writes the specified file into the body stream in an efficient way
// On the client side, the file will typically be downloaded with the given filename
func (c *Context) FileAttachment(filepath, filename string) {
	c.ResW.Header().Set("content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	http.ServeFile(c.ResW, c.ReqW, filepath)
}

// SSEvent writes a Server-Sent Event into the body stream.
// 流式发送数据给客户端
func (c *Context) SSEvent(name string, message interface{}) {
	c.Render(-1, sse.Event{
		Event: name,
		Data:  message,
	})
}

// Stream sends a streaming response and returns a boolean
// indicates "Is client disconnected in middle of stream"
func (c *Context) Stream(step func(w io.Writer) bool) bool {
	w := c.ResW
	clientGone := w.CloseNotify()
	for {
		select {
		case <-clientGone:
			return true
		default:
			keepOpen := step(w)
			w.Flush()
			if !keepOpen {
				return false
			}
		}
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//
/************************************/
/******** CONTENT NEGOTIATION *******/
/************************************/

// Negotiate contains all negotiations data.
type Negotiate struct {
	Offered  []string
	HTMLName string
	HTMLData interface{}
	JSONData interface{}
	XMLData  interface{}
	YAMLData interface{}
	Data     interface{}
}

// Negotiate calls different Render according acceptable Accept format.
func (c *Context) Negotiate(code int, config Negotiate) {
	switch c.NegotiateFormat(config.Offered...) {
	case binding.MIMEJSON:
		data := chooseData(config.JSONData, config.Data)
		c.JSON(code, data)

	case binding.MIMEHTML:
		data := chooseData(config.HTMLData, config.Data)
		c.HTML(code, config.HTMLName, data)

	case binding.MIMEXML:
		data := chooseData(config.XMLData, config.Data)
		c.XML(code, data)

	case binding.MIMEYAML:
		data := chooseData(config.YAMLData, config.Data)
		c.YAML(code, data)

	default:
		c.AbortWithError(http.StatusNotAcceptable, errors.New("the accepted formats are not offered by the server")) // nolint: errcheck
	}
}

// NegotiateFormat returns an acceptable Accept format.
func (c *Context) NegotiateFormat(offered ...string) string {
	assert1(len(offered) > 0, "you must provide at least one offer")

	if c.Accepted == nil {
		c.Accepted = parseAccept(c.requestHeader("Accept"))
	}
	if len(c.Accepted) == 0 {
		return offered[0]
	}
	for _, accepted := range c.Accepted {
		for _, offer := range offered {
			// According to RFC 2616 and RFC 2396, non-ASCII characters are not allowed in headers,
			// therefore we can just iterate over the string without casting it into []rune
			i := 0
			for ; i < len(accepted); i++ {
				if accepted[i] == '*' || offer[i] == '*' {
					return offer
				}
				if accepted[i] != offer[i] {
					break
				}
			}
			if i == len(accepted) {
				return offer
			}
		}
	}
	return ""
}

// SetAccepted sets Accept header data.
func (c *Context) SetAccepted(formats ...string) {
	c.Accepted = formats
}
