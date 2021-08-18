// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

import (
	"errors"
	"time"
)

type sessionKeeper interface {
	Get(string) interface{}
	Set(string, interface{})
	Del(string)
	Save()
	Expire(time.Duration)
}

// GoFast框架的 Context Session
// 默认将使用 Redis 存放 session 信息
type CtxSession struct {
	Sid        string
	Token      string
	TokenIsNew bool
	Saved      bool
	Values     map[string]interface{}
}

// CtxSession 需要实现 sessionKeeper 所有接口
var _ sessionKeeper = &CtxSession{}

// TODO: 你可以自定义实现下面这几个方法，解决底层数据库存储操作。
// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
var CtxSessionSaveFun = func(ss *CtxSession) (string, error) {
	return "", errors.New("Save error. ")
}
var CtxSessionExpireFun = func(ss *CtxSession, ttl time.Duration) (bool, error) {
	return false, errors.New("Change expire error. ")
}
var CtxSessionDestroyFun = func(ss *CtxSession) {}
var CtxSessionCreateFun = func(ctx *Context) {}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (ss *CtxSession) Get(key string) interface{} {
	if ss.Values == nil {
		return nil
	}
	return ss.Values[key]
}

func (ss *CtxSession) Set(key string, val interface{}) {
	ss.Saved = false
	ss.Values[key] = val
}

func (ss *CtxSession) Save() {
	// 如果已经保存了，不会重复保存
	if ss.Saved == true {
		return
	}
	// 调用自定义函数保存当前 session
	_, err := CtxSessionSaveFun(ss)

	// TODO: 如果保存失败怎么办？目前是抛异常，本次请求直接返回错误。
	if err != nil {
		RaisePanic("Save session error.")
	} else {
		ss.Saved = true
	}
}

func (ss *CtxSession) Del(key string) {
	delete(ss.Values, key)
	ss.Saved = false
}

func (ss *CtxSession) Expire(ttl time.Duration) {
	yn, err := CtxSessionExpireFun(ss, ttl)
	if yn == false || err != nil {
		RaisePanic("Session expire error.")
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// 销毁当前Session
func (c *Context) DestroySession() {
	CtxSessionDestroyFun(c.Sess)
	c.Sess = nil
}

func (c *Context) NewSession() {
	CtxSessionCreateFun(c)
}
