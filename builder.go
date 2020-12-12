package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"httpclient/circuitbreaker"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
)

type Builder struct {
	conf Config

	url         string
	method      string
	queryParams *UrlValue
	headers     *Header
	body        []byte
	endpoint    string
	host        string

	ctx      context.Context
	request  *http.Request
	response *http.Response

	handlerIndex int
	handlerChain []HandlerFunc

	filterFunc       func(*http.Request, *http.Response) error
	accessStatusCode []int

	err error

	client *Client
}

func NewBuilder(client *Client) *Builder {
	builder := &Builder{
		client:       client,
		accessStatusCode: []int{http.StatusOK},
		handlerIndex: -1,
		handlerChain: []HandlerFunc{
			GetFilterHandler(),GetBreakerHandler(),
		},
	}

	return builder
}

func (b *Builder) URL(url string) *Builder {
	b.url = url

	return b
}

func (b *Builder) Method(method string) *Builder {
	b.method = method

	return b
}

func (b *Builder) QueryParams(queryParams *UrlValue) *Builder {
	b.queryParams = queryParams

	return b
}

func (b *Builder) Headers(headers *Header) *Builder {
	b.headers = headers
	return b
}

func (b *Builder) Body(body []byte) *Builder {
	b.body = body

	return b
}

func (b *Builder) JsonBody(requestBody interface{}) *Builder {
	if requestBody == nil {
		return b
	}

	requestJson, err := json.Marshal(requestBody)
	if err != nil {
		b.err = err
		return b
	}

	b.body = requestJson

	return b
}

func (b *Builder) FormBody(requestForm *Form) *Builder {
	if requestForm == nil {
		return b
	}

	b.body = []byte(requestForm.Encode())

	return b
}

func (b *Builder) AddHandler(handlerFunc HandlerFunc) *Builder {
	b.handlerChain = append(b.handlerChain, handlerFunc)

	return b
}

func (b *Builder) SetError(err error) *Builder {
	if err != nil {
		b.err = err
	}

	return b
}

func (b *Builder) GetFilterFunc() func(*http.Request, *http.Response) error {
	return b.filterFunc
}

func (b *Builder) SetFilterFunc(filterFunc func(*http.Request, *http.Response) error) *Builder {
	b.filterFunc = filterFunc

	return b
}

func (b *Builder) GetAccessCode() []int {
	return b.accessStatusCode
}

func (b *Builder) SetAccessCode(statusCode ...int) *Builder {
	b.accessStatusCode = statusCode

	return b
}

func (b *Builder) Fetch(ctx context.Context) *Response {
	if b.err != nil {
		NewResponse(nil, b.err)
	}

	if b.queryParams != nil && b.queryParams.Values != nil {
		b.url = fmt.Sprintf("%s?%s", b.url, b.queryParams.Encode())
	}

	if b.headers == nil {
		b.headers = GetDefaultHeader()
	}

	b.ctx = ctx
	b.AddHandler(GetOriginHttpHandler())
	resp, err := b.fetch()

	return NewResponse(resp, err)
}

func (b *Builder) fetch() (*http.Response, error) {
	var bodyReader io.Reader
	if b.body != nil {
		bodyReader = bytes.NewReader(b.body)
	}
	req, err := http.NewRequest(b.method, b.url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header = b.headers.Header

	req = req.WithContext(
		httptrace.WithClientTrace(
			b.ctx,
			&httptrace.ClientTrace{
				GotConn: func(info httptrace.GotConnInfo) {
					b.endpoint = info.Conn.RemoteAddr().String()
				},
			},
		),
	)

	b.host = req.URL.Host
	b.request = req

	b.Next()

	return b.response, b.err
}

func (b *Builder) Next() {
	b.handlerIndex++
	if b.handlerIndex < len(b.handlerChain) {
		b.handlerChain[b.handlerIndex](b)
	}
}

func (b *Builder) DisableBreaker(disableBreaker bool) *Builder {
	b.conf.DisableBreaker = disableBreaker

	return b
}

func (b *Builder) BreakerRate(rate float64) *Builder {
	if rate > 1.0 || rate < 0 {
		b.err = errors.New("breaker rate is invalid")
		return b
	}

	b.conf.BreakerRate = rate
	return b
}

func (b *Builder) BreakerMinSample(minSample int) *Builder {
	if minSample < 0 {
		b.err = errors.New("breaker min sample is invalid")
		return b
	}

	b.conf.BreakerMinSample = minSample
	return b
}

func (b *Builder) GetUrlBreaker(rawUrl string) (*circuitbreaker.Breaker, error) {
	parseUrl, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	if breaker := b.client.breakerGroup.Get(parseUrl.Host); breaker != nil {
		return breaker, nil
	} else {
		breaker := circuitbreaker.NewRateBreaker(b.conf.BreakerRate, int64(b.conf.BreakerMinSample))
		b.client.breakerGroup.Add(parseUrl.Host, breaker)
		return breaker, nil
	}
}
