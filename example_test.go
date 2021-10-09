package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func ExampleNewClient() {
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
		URL("http://www.baidu.com").
		Method(http.MethodGet).
		Headers(GetDefaultHeader()).
		Fetch(context.Background())
	fmt.Println(resp.StatusCode)

	// Output: 200
}
