package httpmux



type treeLeaf struct  {
	name string
	handle MuxHandle
}

type treeNode struct  {
	root *treeLeaf
	leafs map[string]*treeNode
}

type router struct  {
	methods map[string]*method
}

type method struct  {
	staticModeMap map[string]MuxHandle
	dynamicModeMap map[string]*treeNode
}
