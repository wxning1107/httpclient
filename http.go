package httpclient

func NewClient(c *Config) *Client {
	if c == nil {
		panic("http client config is nil")
	}
	if c.RequestTimeout == 0 {
		panic("http client must be set request timeout")
	}
	if c.BreakerMinSample < 0 {
		panic("breaker min sample is invalid")
	}
	if c.BreakerRate > 1.0 || c.BreakerRate < 0 {
		panic("breaker rate is invalid")
	}
	if c.BreakerMinSample == 0 {
		c.BreakerMinSample = 10
	}
	if c.BreakerRate == 0.0 {
		c.BreakerRate = 0.5
	}

	client := Open(c)

	return client
}

