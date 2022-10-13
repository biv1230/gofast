package cst

// 配置http相关设置，比如路由相关控制参数
type WebConfig struct {
	RedirectTrailingSlash bool   `v:"def=false"`     // 探测url后面加减'/'之后是否能匹配路由（这个时代默认不需要了）
	CheckOtherMethodRoute bool   `v:"def=false"`     // 检查其它Method下，是否有对应的路由
	RemoveExtraSlash      bool   `v:"def=false"`     // 规范请求的URL
	UseRawPath            bool   `v:"def=false"`     // 默认取原始的Path，不需要自动转义
	UnescapePathValues    bool   `v:"def=true"`      // 是否把URL中的参数值做转义
	DefNotAllowedHandler  bool   `v:"def=true"`      // 是否采用默认的NotAllowed处理函数
	DefNoRouteHandler     bool   `v:"def=true"`      // 是否采用默认的NoRoute匹配函数
	ForwardedByClientIP   bool   `v:"def=true"`      // 是否从"X-Forwarded-For"的header中提取请求IP地址
	SecureJsonPrefix      string `v:"def=while(1);"` // JsonP安全前缀
	MaxMultipartBytes     int64  `v:"def=33554432"`  // 最大上传文件的大小，默认32MB
	ApplyUrlParamsToPms   bool   `v:"def=true"`      // 将UrlParams解析的参数自动加入Pms
	PrintRouteTrees       bool   `v:"def=false"`     // 是否打印出当前路由数

	//LogType     string `v:"def=json,enum=json|sdx"`              // 日志类型
	//EnableRouteMonitor bool `cnf:",def=true"` // 是否统计路由的访问处理情况，为单个路由的熔断降载做储备
}

// 闪电侠实现的中间件控制参数
type SdxConfig struct {
	// sdx 实现模块的配置参数
	NeedSysCheck     bool   `v:"def=true"`                        // 是否启动CPU使用情况的定时检查工作
	NeedSysPrint     bool   `v:"def=true"`                        // 定时打印系统检查日志
	EnableTimeout    bool   `v:"def=true"`                        // 默认启动超时拦截
	DefTimeoutMS     int64  `v:"def=3000"`                        // 每次请求的超时时间（单位：毫秒）
	MaxContentLength int64  `v:"def=33554432"`                    // 最大请求字节数，32MB（33554432），传0不限制
	MaxConnections   int32  `v:"def=1000000,range=[0:100000000]"` // 最大同时请求数，默认100万同时进入，传0不限制
	JwtSecret        string `v:""`                                // JWT认证的秘钥
}
