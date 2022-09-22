// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

// 强制路由匹配走404逻辑
func (c *Context) SetRouteAs404() {
	c.route.ptrNode = c.myApp.miniNode404
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (c *Context) execHandlers() {
	if c.execIdx == maxRouteHandlers {
		return
	}

	c.RouteIdx = c.route.ptrNode.routeIdx
	c.handlers = c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsItemIdx]
	c.execIdx = -1
	c.Next()
}

// 执行下一个中间件函数
func (c *Context) Next() {
	c.execIdx++
	for c.execIdx < int8(len(c.handlers.hdsIdxChain)) {
		c.myApp.fstMem.tidyHandlers[c.handlers.hdsIdxChain[c.execIdx]](c)
		// 可能被设置成了 abort ，这样后面的 handlers 不用再调用了
		if c.execIdx == maxRouteHandlers {
			break
		}
		c.execIdx++
	}
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (c *Context) execAfterMatchHandlers() {
	if c.route.ptrNode == nil {
		return
	}
	it := c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsItemIdx]
	for it.afterMatchLen > 0 {
		c.myApp.fstMem.tidyHandlers[it.afterMatchIdx](c)
		it.afterMatchLen--
		it.afterMatchIdx++
	}
}

// NOTE: 下面的钩子函数不需要中断执行链。
func (c *Context) execBeforeSendHandlers() {
	if c.route.ptrNode == nil {
		return
	}
	it := c.handlers // c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsItemIdx]
	gp := c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsGroupIdx]

	// 5.preSend
	for it.beforeSendLen > 0 {
		//if c.aborted {
		//	goto over
		//}
		c.myApp.fstMem.tidyHandlers[it.beforeSendIdx](c)
		it.beforeSendLen--
		it.beforeSendIdx++
	}
	for gp.beforeSendLen > 0 {
		//if c.aborted {
		//	goto over
		//}
		c.myApp.fstMem.tidyHandlers[gp.beforeSendIdx](c)
		gp.beforeSendLen--
		gp.beforeSendIdx++
	}
	//over:
	//	return
}

func (c *Context) execAfterSendHandlers() {
	if c.route.ptrNode == nil {
		return
	}
	it := c.handlers // c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsItemIdx]
	gp := c.myApp.fstMem.hdsNodes[c.route.ptrNode.hdsGroupIdx]

	// 6.afterSend
	for it.afterSendLen > 0 {
		//if c.aborted {
		//	goto over
		//}
		c.myApp.fstMem.tidyHandlers[it.afterSendIdx](c)
		it.afterSendLen--
		it.afterSendIdx++
	}
	for gp.afterSendLen > 0 {
		//if c.aborted {
		//	goto over
		//}
		c.myApp.fstMem.tidyHandlers[gp.afterSendIdx](c)
		gp.afterSendLen--
		gp.afterSendIdx++
	}
	//over:
	//	return
}
