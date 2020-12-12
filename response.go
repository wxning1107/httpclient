package httpclient

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
)

type Response struct {
	*http.Response
	err error
}

func newResponse(resp *http.Response, err error) *Response {
	return &Response{
		Response: resp,
		err:      err,
	}
}

func (resp *Response) Error() error {
	if resp.err != nil {
		return resp.err
	} else {
		return nil
	}
}

func (resp *Response) DecodeJson(value interface{}) error {
	if err := resp.Error(); err != nil {
		return err
	}

	if value == nil {
		return errors.New("value is nil")
	}
	if reflect.ValueOf(value).IsNil() {
		return errors.New("reflect value is nil")
	}
	if reflect.TypeOf(value).Kind() != reflect.Ptr {
		return errors.New("value is not ptr")
	}

	defer resp.Response.Body.Close()

	body, err := ioutil.ReadAll(resp.Response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, value)
}

func (resp *Response) Body() (body string, err error) {
	if err := resp.Error(); err != nil {
		return "", err
	}

	defer resp.Response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Response.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

// TODO:
func (resp *Response) IsDegraded() bool {
	return false
}
