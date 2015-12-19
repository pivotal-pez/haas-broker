package fake

import "net/http"

//ClientDoer - a fake implementation of a clientdoer interface
type ClientDoer struct {
	Response   *http.Response
	Error      error
	SpyRequest http.Request
}

//Do - fake do method
func (s *ClientDoer) Do(req *http.Request) (resp *http.Response, err error) {
	s.SpyRequest = *req
	return s.Response, s.Error
}
