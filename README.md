# Elegant Httpclient

### What's this?

The project provides a simple httpclient that is implemented with Golang

### Installation

To install Httpclient package, you need to install Go and set your Go workspace first.

```shell script
$ go get -u github.com/wxning1107/httpclient
```

### How to use?

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
