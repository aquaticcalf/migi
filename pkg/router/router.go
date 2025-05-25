// we love file system based routing, don't we?
package router

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// first let's define each type that we will use in our router

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

// New creates a new Router instance with a given pages directory,
// pages directory -> a directory that contains all the pages of the application, written with gomponents
func New(pagesDir string) (*Router, error) {
	// first create a empty router
	router := &Router{
		routes: make(map[string]*Route),
		tree:   newTree(),
	}
	// now we will walk through the pages directory and register all the routes
	err := filepath.Walk(pagesDir, func(path string, info os.FileInfo, err error) error {
		// we only want to register files that are .go files
		if err != nil || !strings.HasSuffix(path, ".go") {
			return nil
		}
		// extract the relative path from the pages directory
		relativePath, _ := filepath.Rel(pagesDir, path)
		// routePattern is the url path for the route
		routePattern := toPattern(relativePath)
		// route parameters and optional catch all are extracted here
		routeParameters, catch := extractParameters(routePattern)
		// we will create a new route with the pattern and parameters
		// and add it to the router's routes map and the route tree
		router.routes[routePattern] = &Route{
			Pattern:           routePattern,
			Parameters:        routeParameters,
			CatchAllParameter: catch,
		}
		router.tree.addRoute(router.routes[routePattern])
		// continue walking without errors
		return nil
	})
	if err != nil {
		return nil, err
	}

	return router, nil
}

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

// toPattern converts a relative file path to a route
// it removes the ".go" suffix and converts the path to a url like format
// it also handles the special case of "index" files
// for example, "blog/index.go" will be converted to "/blog"
// and "index.go" will be converted to "/"
func toPattern(relativePath string) string {
	urlPath := "/" + strings.TrimSuffix(relativePath, ".go")
	urlPath = filepath.ToSlash(urlPath)
	if strings.HasSuffix(urlPath, "/index") {
		urlPath = strings.TrimSuffix(urlPath, "/index")
		if urlPath == "" {
			urlPath = "/"
		}
	}
	return urlPath
}

// paramRegex is a regular expression to match dynamic parameters and catch all in routes
// for example, it matches "[slug]" in "/blog/[slug]" and "[...parts]" in "/users/[...parts]"
// it captures the parameter name without the brackets or dots
var paramRegex = regexp.MustCompile(`\[\.\.\.([^\]]+)\]|\[([^\]]+)\]`)

// extractParameters extracts dynamic parameters and optional catch all from a route
// it returns a slice of parameter names (excluding catch all dots) and the catch all name or empty
// for example, "/users/[...parts]" returns (["parts"], "parts")
// and "/blog/[slug]" returns (["slug"], "")
// a normal route without parameters like "/about" returns ([], "")
func extractParameters(pattern string) ([]string, string) {
	matches := paramRegex.FindAllStringSubmatch(pattern, -1)
	paramNames := make([]string, 0, len(matches))
	catch := ""
	for _, match := range matches {
		if match[1] != "" {
			// catch all group, because match[1] is not empty
			// match[1] contains the parameter name without the dots
			// for example, "[...parts]" will match and return "parts"
			paramNames = append(paramNames, match[1])
			catch = match[1]
		} else {
			// dynamic single segment parameter, because match[2] is not empty
			// match[2] contains the parameter name without the brackets
			// for example, "[slug]" will match and return "slug"
			paramNames = append(paramNames, match[2])
		}
	}
	return paramNames, catch
}
