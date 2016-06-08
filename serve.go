package httpmux

import (
	"net/http"
	"strings"
)

func (self httpMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m := self.route.methods[r.Method]
	if strings.Contains(r.URL.Path, ":") {
		if handle, get := m.staticModeMap[r.URL.Path]; get {
			c := newContext(r, w, self.res)
			handle(c)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}
	path := newPath(r.URL.Path)
	tree := m.dynamicModeMap
	if handle, get := findHandle(path, tree); get {
		c := newContext(r, w, self.res)
		handle(c)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	return
}

func findHandle(p *path, tree map[string]*treeNode) (MuxHandle, bool) {
	nodeName := p.nextNode()
	if nodeName == "" {
		return nil, false
	}
	node := tree[nodeName]
	if p.subPath == "" {
		h := node.root.handle
		if h != nil {
			return h, true
		}
		return nil, false
	}
	return findHandle(p, node.leafs)
}