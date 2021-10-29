package mid

import (
	"github.com/qinchende/gofast/fst"
	"net/http"
)

// ++++++++++++++++++++++ add by cd.net 2021.10.14
// 总说：定时统计（间隔60秒）系统资源利用情况 | 请求处理相应性能 | 请求量 等
//

func RouteMetric(w *fst.GFResponse, r *http.Request) {
	// 执行请求
	w.NextFit(r)
	w.AddRouteMetric()
}
