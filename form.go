package httpclient

import "net/url"

type Form struct {
	url.Values
}

func NewForm(key, value string) *Form {
	v := url.Values{}
	v.Set(key, value)

	return &Form{
		Values: v,
	}
}

func (f *Form)Add(key, value string) *Form {
	f.Set(key, value)

	return f
}