package httpmux

import (
	"net/http"
	"html/template"
)

func New() *Mux {
	mux := Mux{}
	r := newRouter()
	mux.route = r
	return &mux
}

type Mux struct  {
	route *Router
}

func (m *Mux) Get(path string, handle Handle) {
	m.Handle("GET", path, handle)
}

func (m *Mux) Head(path string, handle Handle) {
	m.Handle("HEAD", path, handle)
}

func (m *Mux) Options(path string, handle Handle) {
	m.Handle("OPTIONS", path, handle)
}

func (m *Mux) post(path string, handle Handle) {
	m.Handle("POST", path, handle)
}

func (m *Mux) Put(path string, handle Handle) {
	m.Handle("PUT", path, handle)
}

func (m *Mux) Patch(path string, handle Handle) {
	m.Handle("PATCH", path, handle)
}

func (m *Mux) Delete(path string, handle Handle) {
	m.Handle("DELETE", path, handle)
}

func (m *Mux) Handle(method, path string, handle Handle) {
	m.route.handle(method, path, handle)
}

func (m *Mux) Template(tpl *template.Template) {
	m.route.tpl = tpl
}

func (m *Mux) Ctx() *context {
	return m.route.ctx
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
func (m *Mux) ServeFiles(path string, root http.FileSystem) {
	m.route.serveFiles(path, root)
}

func (m *Mux) PanicHandlerFunc(handleFunc func(*Request, interface{})) {
	m.route.PanicHandler = handleFunc
}

func (m *Mux) NotFoundHandler(handle Handle) {
	m.route.NotFound = handle
}

func (m *Mux) MethodNotAllowedHandler(handle Handle) {
	m.route.MethodNotAllowed = handle
}

func (m *Mux) RedirectTrailingSlash(flag bool) {
	m.route.RedirectTrailingSlash = flag
}

func (m *Mux) RedirectFixedPath(flag bool) {
	m.route.RedirectFixedPath = flag
}

func (m *Mux) HandleMethodNotAllowed(flag bool) {
	m.route.HandleMethodNotAllowed = flag
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.route.serveHTTP(w, req)
}



