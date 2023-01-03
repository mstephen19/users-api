package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type jsonError struct {
	Message string `json:"message"`
}

func NewJsonError(message string) ([]byte, error) {
	return json.Marshal(jsonError{message})
}

type handler struct {
	// A regular expression for matching the URL's path.
	pattern *regexp.Regexp
	handler http.Handler
}

type router struct {
	// A map of handlers where the keys are request methods and the
	// values are slices of handlers
	handlers map[string][]*handler
	mutex    sync.Mutex
}

// Create a router that can be used to handle routes based on both request method
// and a regular expression that will be used to test the path.
func NewRouter() *router {
	router := &router{}

	router.handlers = make(map[string][]*handler)

	return router
}

func (router *router) Handle(method string, pattern *regexp.Regexp, handle http.Handler) {
	router.mutex.Lock()
	defer router.mutex.Unlock()
	router.handlers[method] = append(router.handlers[method], &handler{pattern, handle})
}

func (router *router) HandleFunc(method string, pattern *regexp.Regexp, handle func(writer http.ResponseWriter, req *http.Request)) {
	router.mutex.Lock()
	defer router.mutex.Unlock()
	router.handlers[method] = append(router.handlers[method], &handler{pattern, http.HandlerFunc(handle)})
}

func (router *router) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	// Remove the "/" from the end of the URL's path to avoid weird cases
	if strings.HasSuffix(req.URL.Path, "/") {
		req.URL.Path = strings.TrimSuffix(req.URL.Path, "/")
	}
	// If the empty path was reached,
	if req.URL.Path == "" {
		req.URL.Path = "/"
	}

	// If the method is an empty string, default to GET
	if req.Method == "" {
		req.Method = http.MethodGet
	}

	if len(router.handlers[req.Method]) == 0 {
		fmt.Println("Zero handlers found for path")
		http.NotFound(writer, req)
		return
	}

	for _, route := range router.handlers[req.Method] {
		if route.pattern.MatchString(req.URL.Path) {
			route.handler.ServeHTTP(writer, req)
			// Only run one matching handler to avoid weird clashes
			return
		}
	}

	http.NotFound(writer, req)
}
