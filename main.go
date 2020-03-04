package main

func main(){
	rr := newPathResolver()
	rr.Add("GET /hello", hello)
	rr.Add("(GET|HEAD) /goodbye(/?[A-Za-z0-9]*)?", goodbye)
}

// resolver
type resolver struct {
	route map[string]string
	cache map[string]*regexp.Regexp
}

// Adding new route
func (r *resolver) Add(regex string) {
	cache, _ := regexp.Compile(regex)
	r.cache[regex] = cache
}

// router
func (r *resolver) router(route string) {
	for pattern := range r.route {
		if r.cache[pattern].MatchString(route) == true {
			// call the function
			return
		}
	}
}