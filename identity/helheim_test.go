// This file was generated by github.com/nelsam/hel.  Do not
// edit this code by hand unless you *really* know what you're
// doing.  Expect any changes made manually to be overwritten
// the next time hel regenerates this file.

package identity_test

import (
	"net/http"
)

type mockHTTPClient struct {
	GetCalled chan bool
	GetInput  struct {
		Url chan string
	}
	GetOutput struct {
		Resp chan *http.Response
		Err  chan error
	}
}

func newMockHTTPClient() *mockHTTPClient {
	m := &mockHTTPClient{}
	m.GetCalled = make(chan bool, 100)
	m.GetInput.Url = make(chan string, 100)
	m.GetOutput.Resp = make(chan *http.Response, 100)
	m.GetOutput.Err = make(chan error, 100)
	return m
}
func (m *mockHTTPClient) Get(url string) (resp *http.Response, err error) {
	m.GetCalled <- true
	m.GetInput.Url <- url
	return <-m.GetOutput.Resp, <-m.GetOutput.Err
}
