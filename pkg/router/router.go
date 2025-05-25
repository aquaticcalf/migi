// we love file system based routing, don't we?
package router

import (
	"os"
	"path/filepath"
	"strings"
)

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
