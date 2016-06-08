package httpmux

import "strings"

type path struct  {
	path string
	subPath string
}

func newPath(path string) *path {
	return path{path:strings.TrimSpace(path), subPath:""}
}

func (self *path) nextNode() string  {
	if self.path == "" {
		return ""
	}
	p := self.path
	idx := strings.Index(p, "/")
	if idx < 0 {
		self.subPath = ""
		return p
	}
	if idx == 0 {
		p = p[1:]
	}
	idx = strings.Index(p, "/")
	if idx < 0 {
		self.subPath = ""
		return p
	}
	self.subPath = p[idx:]
	return p[0:idx]
}

func (self *path) patch(pathTree map[string]*treeNode, handle MuxHandle) map[string]*treeNode {
	nodeName := self.nextNode()
	if nodeName == "" {
		return pathTree
	}
	var node *treeNode
	var nodeExist bool
	if node, nodeExist = pathTree[nodeName]; !nodeExist {
		node = &treeNode{}
		node.root = &treeLeaf{name:nodeName}
		node.leafs = make(map[string]*treeNode)
		pathTree[nodeName] = node
	}
	if self.subPath == "" {
		node.root.handle = handle
		return pathTree
	}
	node.leafs = self.patch(node.leafs, handle)
	return pathTree
}