// Copyright 2022 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package sdx

import "github.com/qinchende/gofast/fst"

func PmsParser(c *fst.Context) {
	if err := c.BuildPms(); err != nil {
		c.AbortFai(0, "PmsParser error: "+err.Error())
	}
}

//func JwtAuthHandler(secret string) fst.CtxHandler {
//	return mid.JwtAuthHandler(secret)
//}
//func BuildPmsOfJson(ctx *fst.Context) {
//	ctx.GenPmsByJSONBody()
//}
//
//func BuildPmsOfXml(ctx *fst.Context) {
//	ctx.GenPmsByXMLBody()
//}
//
//func BuildPmsOfForm(ctx *fst.Context) {
//	ctx.GenPmsByFormBody()
//}
