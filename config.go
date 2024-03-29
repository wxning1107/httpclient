package httpclient

import "time"

const (
	DefaultMaxIdleConns        = 100
	DefaultMaxIdleConnsPerHost = 2
	DefaultIdleConnTimeout     = time.Second * 90
)

type Config struct {
	// 	请求超时时间
	RequestTimeout time.Duration `json:"request_timeout" yaml:"request_timeout"`

	// 是否保持长连接
	DisableKeepAlives bool `json:"disable_keep_alives" yaml:"disable_keep_alives"`
	// 最大空闲连接
	MaxIdleConns int `json:"max_idle_conns" yaml:"max_idle_conns"`
	// 每个host的最大空闲连接
	MaxIdleConnsPerHost int `json:"max_idle_conns_per_host" yaml:"max_idle_conns_per_host"`
	// 每个host的最大连接
	MaxConnsPerHost int `json:"max_conns_per_host" yaml:"max_conns_per_host"`
	// 空闲连接超时时间
	IdleConnTimeout time.Duration `json:"idle_conn_timeout" yaml:"idle_conn_timeout"`

	// 是否关闭断路器
	DisableBreaker bool `json:"disable_breaker" yaml:"disable_breaker"`
	// 断路器断路最小错误比例 [0,1]
	BreakerRate float64 `json:"breaker_rate" yaml:"breaker_rate"`
	// 断路器断路最小采样数
	BreakerMinSample int `json:"breaker_min_sample" yaml:"breaker_min_sample"`
}
