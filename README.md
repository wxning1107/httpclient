# Element Httpclient

### What's this?

The project provides a simple httpclient that is implemented with Golang

## How to use?

```go
	client := NewClient(&Config{
		RequestTimeout: time.Second,
	})

	resp := client.Builder().
		URL("http://www.google.com").
		Method(http.MethodGet).
		Headers(GetDefaultHeader()).
		Fetch(context.Background())
```
