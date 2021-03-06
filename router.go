package httpmux

import (
	"net/http"
	"html/template"
)

// Handle is a function that can be registered to a route to handle HTTP
// requests. Like http.HandlerFunc, but has a third parameter for the values of
// wildcards (variables).
//type Handle func(http.ResponseWriter, *http.Request, Params)

type Handle func(*Request)

func newContext() *context {
	ctx := context{
		data : make(map[string]interface{}),
	}
	return &ctx
}

type context struct {
	data map[string]interface{}
}

func (ctx *context) Set(name string, val interface{})  {
	ctx.data[name] = val
}

func (ctx *context) Get(name string) (bool, interface{}) {
	if val, ok := ctx.data[name]; ok {
		return true, val
	}
	return false, nil
}

func (ctx *context) MustGet(name string) interface{} {
	if val, ok := ctx.data[name]; ok {
		return val
	}
	panic("ctx miss " + name)
	return nil
}

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {
	trees map[string]*node

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// Configurable http.HandlerFunc which is called when no matching route is
	// found. If it is not set, http.NotFound is used.
	NotFound Handle

	// Configurable http.HandlerFunc which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	MethodNotAllowed Handle

	// Function to handle panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler func(*Request, interface{})

	tpl *template.Template

	ctx *context
}


// New returns a new initialized Router.
// Path auto-correction, including trailing slashes, is enabled by default.
func newRouter() *Router {
	return &Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		ctx : newContext(),
	}
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Router) handle(method, path string, handle Handle) {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if r.trees == nil {
		r.trees = make(map[string]*node)
	}

	root := r.trees[method]
	if root == nil {
		root = new(node)
		r.trees[method] = root
	}

	root.addRoute(path, handle)
}


// ServeFiles serves files from the given file system root.
// The path must end with "/*filepath", files are then served from the local
// path /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use http.Dir:
//     router.ServeFiles("/src/*filepath", http.Dir("/var/www"))
func (r *Router) serveFiles(path string, root http.FileSystem) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path '" + path + "'")
	}
	fileServer := http.FileServer(root)
	r.handle("GET", path, func(r *Request) {
		req := r.r
		req.URL.Path = r.Query("filepath")
		fileServer.ServeHTTP(r.w.w, req)
	})
}

func (r *Router) recv(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		request := newRequest(w, req, r.tpl, r.ctx)
		r.PanicHandler(request, rcv)
	}
}

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handle function and the path parameter
// values. Otherwise the third return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *Router) Lookup(method, path string) (Handle, Params, bool) {
	if root := r.trees[method]; root != nil {
		return root.getValue(path)
	}
	return nil, nil, false
}

// ServeHTTP makes the router implement the http.Handler interface.
func (r *Router) serveHTTP(w http.ResponseWriter, req *http.Request) {
	if r.PanicHandler != nil {
		defer r.recv(w, req)
	}
	if root := r.trees[req.Method]; root != nil {
		path := req.URL.Path

		if handle, ps, tsr := root.getValue(path); handle != nil {
			request := newRequest(w, req, r.tpl, r.ctx)
			request.params = ps
			handle(request)
			return
		} else if req.Method != "CONNECT" && path != "/" {
			code := 301 // Permanent redirect, request with GET method
			if req.Method != "GET" {
				// Temporary redirect, request with same method
				// As of Go 1.3, Go does not support status code 308.
				code = 307
			}

			if tsr && r.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}
				http.Redirect(w, req, req.URL.String(), code)
				return
			}

			// Try to fix the request path
			if r.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					CleanPath(path),
					r.RedirectTrailingSlash,
				)
				if found {
					req.URL.Path = string(fixedPath)
					http.Redirect(w, req, req.URL.String(), code)
					return
				}
			}
		}
	}

	// Handle 405
	if r.HandleMethodNotAllowed {
		for method := range r.trees {
			// Skip the requested method - we already tried this one
			if method == req.Method {
				continue
			}

			handle, _, _ := r.trees[method].getValue(req.URL.Path)
			if handle != nil {
				if r.MethodNotAllowed != nil {
					request := newRequest(w, req, r.tpl, r.ctx)
					request.w.Status(http.StatusMethodNotAllowed)
					r.MethodNotAllowed(request)
				} else {
					http.Error(w,
						http.StatusText(http.StatusMethodNotAllowed),
						http.StatusMethodNotAllowed,
					)
				}
				return
			}
		}
	}

	// Handle 404
	if r.NotFound != nil {
		request := newRequest(w, req, r.tpl, r.ctx)
		request.w.Status(http.StatusNotFound)
		r.NotFound(request)
	} else {
		http.NotFound(w, req)
	}
}




