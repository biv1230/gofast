package fstx

import (
	"github.com/qinchende/gofast/fst"
	"github.com/qinchende/gofast/fst/gate"
	"github.com/qinchende/gofast/fst/mid"
	"time"
)

// 第一级：
// NOTE：Fit系列是全局的，针对所有请求起作用，而且不区分路由，这个时候还根本没有开始匹配路由。
// GoFast提供默认的全套拦截器，开启微服务治理
// 请求按照先后顺序依次执行这些拦截器，顺序不可随意改变
func DefaultFits(gft *fst.GoFast) *fst.GoFast {
	gft.Fit(mid.MaxConnections(gft.FitMaxConnections)) // 最大同时处理请求数量：100万
	//gft.Fit(mid.CpuMetric(nil))                        // cpu 统计 | 熔断

	return gft
}

// 第二级：
// 带上下文 gofast.fst.Context 的执行链
func DefaultHandlers(gft *fst.GoFast) *fst.GoFast {
	// 初始化一个全局的 请求管理器（记录访问数据，分析统计，限流熔断，定时日志）
	reqKeeper := gate.CreateReqKeeper(gft.AppName(), gft.FullPath)
	// 因为Routes的数量只能在加载完所有路由之后才知道
	// 所以这里选择延时构造所有Breakers
	gft.OnBeforeBuildRoutes(func(gft *fst.GoFast) {
		reqKeeper.SetBreakers(gft.RouteLength())
	})

	gft.Before(mid.Tracing)                                                         // 链路追踪，在日志打印之前执行，日志才会体现出标记
	gft.Before(mid.ReqLogger)                                                       // 所有请求写日志，根据配置输出日志样式
	gft.Before(mid.HardwareMetric(nil))                                             // 硬件资源利用率高，答应信息并熔断
	gft.Before(mid.Breaker(reqKeeper))                                              // 针对不同路由的配置，启动熔断机制
	gft.Before(mid.ReqTimeout(time.Duration(gft.FitReqTimeout) * time.Millisecond)) // 超时自动返回，后台处理继续，默认3000毫秒
	gft.Before(mid.Recovery)                                                        // 截获所有异常
	gft.Before(mid.TimeMetric(reqKeeper))                                           // 请求耗时统计
	gft.Before(mid.Prometheus)                                                      // 适合 prometheus 的统计信息
	gft.Before(mid.MaxContentLength(gft.FitMaxContentLength))                       // 最大的请求头限制，默认32MB
	gft.Before(mid.Gunzip)                                                          // 自动 gunzip 解压缩

	// 下面的这些特性恐怕都需要用到 fork 时间模式添加监控。
	// gft.Fit(mid.JwtAuthorize(gft.FitJwtSecret))
	return gft
}

// ++++++++++++++++++++++++++++++++++ go-zero default handler chains
//  chain := alice.New(
//  handler.TracingHandler,
//  s.getLogHandler(),
//  handler.MaxConns(s.conf.MaxConns),
//  handler.BreakerHandler(route.Method, route.Path, metrics),
//  handler.SheddingHandler(s.getShedder(fr.priority), metrics),
//  handler.TimeoutHandler(time.Duration(s.conf.Timeout)*time.Millisecond),
//  handler.RecoverHandler,
//  handler.MetricHandler(metrics),
//  handler.PromethousHandler(route.Path),
//  handler.MaxBytesHandler(s.conf.MaxBytes),
//  handler.GunzipHandler,
//  )
//  chain = s.appendAuthHandler(fr, chain, verifier)
//
//  for _, middleware := range s.middlewares {
//  chain = chain.Append(convertMiddleware(middleware))
//  }
//  handle := chain.ThenFunc(route.Handler)
