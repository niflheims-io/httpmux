package httpmux

import (
	"strings"
)

func New() *httpMux {
	mux := httpMux{}
	route := router{}
	route.methods = make(map[string]*method)
	mux.route = route
	mux.res = new(resMap)
	return &mux
}

type httpMux struct  {
	route *router
	res *resMap
}

type MuxHandle func(*Context)

type resMap map[string]interface{}

func (self *httpMux) pathMap(method, path string, handle MuxHandle)  {
	if strings.Index(path, "/") != 0 {
		path = "/" + path
	}
	var methods *method
	var methodExist bool
	if methods, methodExist = self.route.methods[method]; !methodExist {
		methods = &method{}
		methods.dynamicModeMap = make(map[string]*treeNode)
		methods.staticModeMap = make(map[string]MuxHandle)
		self.route["GET"] = methods
	}
	static := strings.Contains(path, ":")
	if static {
		methods.staticModeMap[path] = &handle
		return
	}
	methods.dynamicModeMap = newPath(path).patch(methods.dynamicModeMap, handle)
}

func (self *httpMux) Get(path string, handle MuxHandle)  {
	self.pathMap("GET", path, &handle)
}

func (self *httpMux) Post(path string, handle MuxHandle)  {
	self.pathMap("POST", path, &handle)
}

func (self *httpMux) Put(path string, handle MuxHandle)  {
	self.pathMap("PUT", path, &handle)
}

func (self *httpMux) Delete(path string, handle MuxHandle)  {
	self.pathMap("DELETE", path, &handle)
}

func (self *httpMux) Trace(path string, handle MuxHandle)  {
	self.pathMap("TRACE", path, &handle)
}

func (self *httpMux) Head(path string, handle MuxHandle)  {
	self.pathMap("HEAD", path, &handle)
}

func (self *httpMux) Options(path string, handle MuxHandle)  {
	self.pathMap("OPTIONS", path, &handle)
}


