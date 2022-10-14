// Copyright 2020 GoFast Author(http://chende.ren). All rights reserved.
// Use of this source code is governed by a MIT license
package fst

func (gft *GoFast) RouteLen() uint16 {
	return uint16(len(gft.allRoutes))
}

func (gft *GoFast) RoutePaths() []string {
	allPaths := make([]string, len(gft.allRoutes))
	for i := 0; i < len(allPaths); i++ {
		allPaths[i] = gft.allRoutes[i].fullPath
	}
	return allPaths
}

// 获取相应路由节点完整URL
func (gft *GoFast) FullPath(idx uint16) string {
	if idx < 0 || int(idx) >= len(gft.allPaths) {
		return ""
	}
	return gft.allPaths[idx]
}

func (c *Context) FullPath() string {
	if c.route.ptrNode != nil {
		return c.myApp.allPaths[c.route.ptrNode.routeIdx]
	} else {
		return ""
	}
}

//
//func (ri *RouteItem) FullPath() string {
//	return ri.fullPath
//}

//func (c *Context) CurrRoute() *RouteItem {
//	if c.RouteIdx <= 0 || c.RouteIdx >= uint16(len(c.myApp.allRoutes)) {
//		return nil
//	}
//	return c.myApp.allRoutes[c.RouteIdx]
//}
