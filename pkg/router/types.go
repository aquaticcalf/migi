package router

// Router manages file based routing
type Router struct {
	routes map[string]*Route
	tree   *routeTree
}

// Route represents a single route in the router
type Route struct {
	Pattern           string   // route pattern, example : "/pictures/1" or "/blog/[slug]" or "/users/[...parts]"
	Parameters        []string // parameter names in order of appearance, example : "/user/[username]/post/[postId]" will have Parameters ["username", "postId"]
	CatchAllParameter string   // name of the catch all parameter if used, example : "/users/[...parts]" will have CatchAllParameter "parts", it is empty by default
}

// RouteTree is a tree structure for efficient route matching
type routeTree struct {
	root *routeTreeNode
}

// routeTreeNode represents a node in the route tree
type routeTreeNode struct {
	children      map[string]*routeTreeNode // static children for segments, example : "blog" in "/blog/[slug]"
	dynamicChild  *routeTreeNode            // dynamic child for parameters, example : "/blog/[slug]" will have a dynamic child node for "[slug]"
	catchAllChild *routeTreeNode            // catch all child for multi segment params, example : "/users/[...parts]"
	parameter     string                    // parameter name for dynamic or catch all routes, example : "slug" or "parts"
	route         *Route                    // associated route for this node, example : "/blog/[slug]" will have a route with Pattern "/blog/[slug]" and Parameters ["slug"]
}
