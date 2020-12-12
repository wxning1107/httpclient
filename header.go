package httpclient

import "net/http"

type Header struct {
	http.Header
}

func GetDefaultHeader() *Header {
	return NewJsonHeader()
}

func NewJsonHeader() *Header {
	header := http.Header{}
	header.Set("Content-Type", "application/json")

	return &Header{
		Header: header,
	}
}

func NewFormURLEncodedHeader() *Header {
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")

	return &Header{
		Header: header,
	}
}

func (h *Header) Add(key, value string) *Header {
	h.Set(key, value)

	return h
}

