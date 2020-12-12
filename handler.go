package httpclient

import (
	"errors"
)

type HandlerFunc func(*Builder)

func GetOriginHttpHandler() HandlerFunc {
	return func(b *Builder) {
		resp, err := b.client.client.Do(b.request)
		b.response = resp
		b.SetError(err)
	}
}

// TODO:
func GetBreakerHandler() HandlerFunc {
	return func(b *Builder) {
		if b.conf.DisableBreaker {
			b.next()
			return
		}

		//breaker, err := b.GetUrlBreaker(b.url)
		//if err != nil {
		//	b.SetError(err)
		//	return
		//}

		//breaker.Call()
	}
}

func GetFilterHandler() HandlerFunc {
	return func(b *Builder) {
		b.next()

		filterFunc := b.GetFilterFunc()
		if filterFunc != nil {
			err := filterFunc(b.request, b.response)
			b.SetError(err)
		} else {
			GetAccessStatusCodeHandler()(b)
		}
	}
}

func GetAccessStatusCodeHandler() HandlerFunc {
	return func(b *Builder) {
		b.next()

		resp := b.response
		if resp == nil {
			return
		}

		isAccess := len(b.GetAccessCode()) == 0
		for _, code := range b.GetAccessCode() {
			if resp.StatusCode == code {
				isAccess = true
				break
			}
		}
		if !isAccess {
			b.SetError(errors.New("wrong status code"))
		}
	}
}
