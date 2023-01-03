package lib

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Handler struct {
	// A regular expression for matching the URL's path.
	pattern *regexp.Regexp
	handler http.Handler
}

type Router struct {
	// A map of handlers where the keys are request methods and the
	// values are slices of handlers
	handlers map[string][]*Handler
}

// Create a router that can be used to handle routes based on both request method
// and a regular expression that will be used to test the path.
func NewRouter() *Router {
	router := &Router{}
	router.handlers = make(map[string][]*Handler)
	return router
}

func (router *Router) Handle(method string, pattern *regexp.Regexp, handler http.Handler) {
	router.handlers[method] = append(router.handlers[method], &Handler{pattern, handler})
}

func (router *Router) HandleFunc(method string, pattern *regexp.Regexp, handler func(writer http.ResponseWriter, req *http.Request)) {
	router.handlers[method] = append(router.handlers[method], &Handler{pattern, http.HandlerFunc(handler)})
}

func (router *Router) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	// Remove the "/" from the end of the URL's path to avoid weird cases
	if strings.HasSuffix(req.URL.Path, "/") {
		req.URL.Path = strings.TrimSuffix(req.URL.Path, "/")
	}

	// If the empty path was reached,
	if req.URL.Path == "" {
		req.URL.Path = "/"
	}

	if len(router.handlers[req.Method]) == 0 {
		fmt.Println("Zero handlers found for path")
		http.NotFound(writer, req)
		return
	}

	for _, route := range router.handlers[req.Method] {
		if route.pattern.MatchString(req.URL.Path) {
			func() { route.handler.ServeHTTP(writer, req) }()
			// Only run one matching handler to avoid weird clashes
			return
		}
	}

	http.NotFound(writer, req)
}
