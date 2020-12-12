package httpclient

import (
	"context"
	"httpclient/circuitbreaker"
	"net"
	"net/http"
	"time"
)

type Client struct {
	// 实际客户端
	client *http.Client
	// 断路器组
	breakerGroup *circuitbreaker.BreakerGroup
	// 全局上下文
	globalContext context.Context
	// 全局上下文取消函数
	globalCancelFunc context.CancelFunc
}

func open(c *Config) *Client {
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = DefaultMaxIdleConns
	}
	if c.MaxIdleConnsPerHost == 0 {
		c.MaxIdleConnsPerHost = DefaultMaxIdleConnsPerHost
	}
	if c.IdleConnTimeout == 0 {
		c.IdleConnTimeout = DefaultIdleConnTimeout
	}

	client := new(Client)
	client.globalContext, client.globalCancelFunc = context.WithCancel(context.Background())
	client.breakerGroup = circuitbreaker.NewBreakerGroup()
	client.client = &http.Client{
		Timeout: c.RequestTimeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			DisableKeepAlives:     c.DisableKeepAlives,
			MaxIdleConns:          c.MaxIdleConns,
			MaxIdleConnsPerHost:   c.MaxIdleConnsPerHost,
			MaxConnsPerHost:       c.MaxConnsPerHost,
			IdleConnTimeout:       c.IdleConnTimeout,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: time.Second,
		},
	}

	return client
}

func (c *Client) Builder() *Builder {
	builder := newBuilder(c)

	return builder
}

func (c *Client) Close() error {
	c.globalCancelFunc()

	return nil
}
