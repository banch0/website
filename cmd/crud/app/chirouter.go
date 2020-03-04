package app

import "net/http"

func chiRouter() {

}

// ChiMux ...
type ChiMux struct {
	MethodNotAllowed func(h http.HandlerFunc)
	NotFound         func(h http.HandlerFunc)
}

func (router *ChiMux) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	router.ServeHTTP(res, req)
}
