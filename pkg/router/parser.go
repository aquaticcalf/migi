package router

import (
	"path/filepath"
	"regexp"
	"strings"
)

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
