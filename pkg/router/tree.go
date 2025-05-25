package router

import "strings"

// newTree initializes a new route tree
func newTree() *routeTree {
	return &routeTree{
		root: &routeTreeNode{children: make(map[string]*routeTreeNode)},
	}
}

// addRoute adds a route to the route tree
func (t *routeTree) addRoute(route *Route) {
	// split the pattern into segments
	pathSegments := strings.Split(strings.Trim(route.Pattern, "/"), "/")
	currentNode := t.root
	// iterate through the segments and add them to the tree
	for _, segment := range pathSegments {
		// catch all parameter, matches any number of segments
		if strings.HasPrefix(segment, "[...") {
			paramName := segment[4 : len(segment)-1]
			if currentNode.catchAllChild == nil {
				currentNode.catchAllChild = &routeTreeNode{
					children:  make(map[string]*routeTreeNode),
					parameter: paramName,
				}
			}
			currentNode = currentNode.catchAllChild
			break // catch all must be last

		} else if strings.HasPrefix(segment, "[") { // dynamic single segment parameter
			paramName := segment[1 : len(segment)-1]
			if currentNode.dynamicChild == nil {
				currentNode.dynamicChild = &routeTreeNode{
					children:  make(map[string]*routeTreeNode),
					parameter: paramName,
				}
			}
			currentNode = currentNode.dynamicChild

		} else { // static segment
			if currentNode.children[segment] == nil {
				currentNode.children[segment] = &routeTreeNode{children: make(map[string]*routeTreeNode)}
			}
			currentNode = currentNode.children[segment]
		}
	}
	// finally, we set the route for the current node
	currentNode.route = route
}
