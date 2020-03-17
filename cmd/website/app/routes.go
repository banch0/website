package app

import (
	"net/http"
	"regexp"
)

// за разделение handler'ов по адресам -> routing
func (receiver *server) InitRoutes() {
	mux := receiver.router.(*ExactMux)
	// panic, если происходят конфликты
	// Handle - добавляет Handler (неудобно)
	// HandleFunc

	// стандартный mux:
	// - если адрес начинается со "/" - то под действие обработчика попадает всё, что начинается со "/"
	mux.GET("/", receiver.handleBurgersList())

	//mux.POST("/burgers/save", receiver.handleBurgersSave())
	//mux.POST("/burgers/remove", receiver.handleBurgersRemove())

	// - но если есть более "специфичный", то используется он
	mux.GET("/favicon.ico", receiver.handleFavicon())
}

// not have duplicate controll, instead use last object

func (receiver *server) InitRoutesPath() {
	mux := receiver.router.(*Resolver)
	mux.Add("GET /hello", hello)
	mux.Add("(GET|HEAD) /goodbye(/?[A-Za-z0-9]*)?", goodbye)
	mux.Add("(GET|HEAD) /goodbye(/?[A-Za-z0-9]*)?", good)
	mux.Add("(GET|HEAD|POST) /test(/?[A-Za-z0-9]*)?", test)
}

func test(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("test \n"))
}

func hello(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("hello"))
}
func goodbye(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("goodbye"))
}
func good(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("goodbye2"))
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
