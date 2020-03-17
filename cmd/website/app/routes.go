package app

import (
	"net/http"
	"regexp"
)

func (receiver *server) InitRoutesPath() {
	mux := receiver.router.(*Resolver)
	mux.Add("(GET|HEAD) /goodbye(/?[A-Za-z0-9]*)?", goodbye)
	mux.Add("GET /api/burgers", receiver.handleAllBurgers())
	mux.Add("GET /api/burger(/?[0-9]*)?", receiver.getBurgerByID())
	mux.Add("(POST|HEAD) /api/burgers", receiver.handleBurgerSave())
	mux.Add("(PUT|HEAD) /api/burgers", receiver.handleBurgerUpdate())
	mux.Add("(DELETE|HEAD) /api/burgers(/?[0-9]*)?", receiver.handleDeleteBurgers())
}

func goodbye(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("goodbye"))
}

// Resolver ...
type Resolver struct {
	handlers map[string]http.HandlerFunc
	cache    map[string]*regexp.Regexp
}

// NewPathResolver ...
func NewPathResolver() *Resolver {
	return &Resolver{
		handlers: make(map[string]http.HandlerFunc),
		cache:    make(map[string]*regexp.Regexp),
	}
}

// Add new route
func (r *Resolver) Add(regex string, handler http.HandlerFunc) {
	r.handlers[regex] = handler
	cache, _ := regexp.Compile(regex)
	r.cache[regex] = cache
}

// ServeHTTP ...
func (r *Resolver) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	check := req.Method + " " + req.URL.Path
	for pattern, handlerFunc := range r.handlers {
		if r.cache[pattern].MatchString(check) == true {
			handlerFunc(res, req)
			return
		}
	}
	http.NotFound(res, req)
}
