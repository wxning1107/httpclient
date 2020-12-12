package httpclient

import "net/url"

type UrlValue struct {
	url.Values
}

func NewQueryParams(key, value string) *UrlValue {
	v := url.Values{}
	v.Set(key, value)

	return &UrlValue{
		Values: v,
	}
}

func (v *UrlValue) Add(key, value string) *UrlValue {
	v.Values.Add(key, value)

	return v
}
