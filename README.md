# Element Httpclient

### What's this?

The project provides a simple httpclient that is implemented with Golang

## How to use?

```go
	client := NewClient(&Config{
		RequestTimeout:      time.Second,
		DisableKeepAlives:   false,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 2,
		MaxConnsPerHost:     0,
		IdleConnTimeout:     time.Second * 90,
		DisableBreaker:      true,
	})

	resp := client.Builder().
		URL("http://www.google.com").
		Method("GET").
		Headers(GetDefaultHeader()).
		Fetch(context.Background())
```
