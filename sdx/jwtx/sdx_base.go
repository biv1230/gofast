package jwtx

import (
	"github.com/qinchende/gofast/connx/gfrds"
)

var (
	sdxTokenPrefix   = "t:"
	sdxSessKeyPrefix = "tls:"
)

type SdxSessCnf struct {
	RedisConnCnf gfrds.ConnCnf `v:"required"`                       // 用 Redis 做持久化
	CheckTokenIP bool          `v:"def=true"`                       // 看是否检查 token ip 地址
	AuthField    string        `v:"def=user_id"`                    // 标记当前登录用户字段是 user_id
	Secret       string        `v:"required=true"`                  // token秘钥
	TTL          int32         `v:"def=14400,range=[0:2000000000]"` // session有效期 默认 3600*4 秒
	TTLNew       int32         `v:"def=180,range=[0:2000000000]"`   // 首次产生的session有效期 默认 60*3 秒
}

// 参数配置，Redis实例等
type SdxSession struct {
	SdxSessCnf
	Redis   *gfrds.GfRedis
	isReady bool // 是否已经初始化
}

// 每个进程只有一个全局 SdxSS 配置对象
var SdxSS *SdxSession