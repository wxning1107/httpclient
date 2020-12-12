package httpclient

import (
	"time"
)

func ExampleHttpclient() {
	client := NewClient(&Config{
		RequestTimeout: time.Second,
	})

	err := client.Builder().URL("www.baidu.com").Method("Get").Headers(GetDefaultHeader()).Fetch()
	if err != nil {
		panic(err)
	}

	// Output: ""
}
