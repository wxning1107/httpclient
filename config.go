package httpclient

import "time"

const (
	DefaultMaxIdleConns        = 100
	DefaultMaxIdleConnsPerHost = 2
	DefaultIdleConnTimeout     = time.Second * 90
)

type Config struct {
	// 	请求超时时间
	RequestTimeout time.Duration

	// 是否保持长连接
	DisableKeepAlives bool
	// 最大空闲连接
	MaxIdleConns int
	// 每个host的最大空闲连接
	MaxIdleConnsPerHost int
	// 每个host的最大连接
	MaxConnsPerHost int
	// 空闲连接超时时间
	IdleConnTimeout time.Duration

	// 是否关闭断路器
	DisableBreaker bool
	// 断路器断路最小错误比例 [0,1]
	BreakerRate float64
	// 断路器断路最小采样数
	BreakerMinSample int
}
