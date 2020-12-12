package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func ExampleHttpclient() {
	client := NewClient(&Config{
		RequestTimeout: time.Second,
	})

	resp := client.Builder().
		URL("http://www.baidu.com").
		Method(http.MethodGet).
		Headers(GetDefaultHeader()).
		Fetch(context.Background())
	fmt.Println(resp.StatusCode)

	// Output: 200
}
