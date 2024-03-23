package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"qrdb/qrdb/pkg/server/myctx"
	"qrdb/qrdb/pkg/server/routes/qrcode"
	"strings"
)

type Routable interface {
	ServeHTTP(r *http.Request, ctx context.Context) (int, interface{})
}

func (s *server) RegisterRoute(
	path string,
	routable Routable,
	allowedMethods []string,
	allowedRoles []string,
) {
	s.pathToRoutable[path] = routable
	s.pathToAllowedMethods[path] = allowedMethods
	s.pathToAllowedRoles[path] = allowedRoles
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	log.Printf("%s %s", r.Method, r.URL.Path)

	requestRoles := []string{"role1", "Y"}

	ctx = context.WithValue(ctx, myctx.Roles("roles"), requestRoles)

	// Find a routable from the path
	if routable, params, route := s.findRoute(r.URL.Path); routable != nil {

		for keyParam, valueParam := range params {
			ctx = context.WithValue(ctx, myctx.Param(keyParam), valueParam)
		}

		// check if the method is allowed
		allowedMethods := s.pathToAllowedMethods[route]
		methodAllowed := len(allowedMethods) == 0
		for _, allowedMethod := range allowedMethods {
			if allowedMethod == r.Method {
				methodAllowed = true
				break
			}
		}
		if !methodAllowed {
			sendReply(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		// Check if it is protected by roles
		routableRoles := s.pathToAllowedRoles[route]
		authorized := len(routableRoles) == 0
		for _, routableRole := range routableRoles {
			if authorized {
				break
			}
			for _, requestRole := range requestRoles {
				if routableRole == requestRole {
					authorized = true
					break
				}
			}
		}
		if !authorized {
			sendReply(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		httpCode, output := routable.ServeHTTP(r, ctx)
		sendReply(w, httpCode, output)

		return
	}

	sendReply(w, http.StatusNotFound, "not found")
}

func sendReply(w http.ResponseWriter, httpCode int, output interface{}) {
	if httpCode >= 400 {
		fmt.Print(output)
		if outputString, ok := output.(string); ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(httpCode)
			fmt.Fprintf(w, `{"success": false, "error": "%s"}`, outputString)
			return
		}
	}

	// Handle raw bytes
	if binaryData, ok := output.(qrcode.RawPNG); ok {
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(httpCode)
		w.Header().Set("Content-Type", "image/png")
		// Scrive i dati binari grezzi come risposta HTTP
		w.Write(binaryData)

		return
	}

	// Try to deserialize the reply
	bytes, err := json.Marshal(output)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpCode)
		fmt.Fprint(w, string(bytes))
		return
	}

	w.WriteHeader(httpCode)
	fmt.Fprint(w, output)
}

func (s *server) findRoute(path string) (Routable, map[string]string, string) {
	params := make(map[string]string)

	// Remove trailing "/"
	path = strings.TrimRight(path, "/")

	// Exact Route
	if route, exists := s.pathToRoutable[path]; exists {
		return route, params, path
	}

	// Extract params
	pathParts := getParts(path)

	// Do not allow parts starting with ":"
	for _, part := range pathParts {
		if strings.HasPrefix(part, ":") {
			return nil, params, ""
		}
	}

	// Do not allow parts starting with "."
	for _, part := range pathParts {
		if strings.HasPrefix(part, ".") {
			return nil, params, ""
		}
	}

	// Check Route by route
	for route, routable := range s.pathToRoutable {
		routeParts := getParts(route)
		if len(routeParts) != len(pathParts) {
			continue
		}

		found := true
		pathParams := make(map[string]string)
		for i, pathPart := range pathParts {
			routePart := routeParts[i]
			if routePart == pathPart {
				continue
			} else if strings.HasPrefix(routePart, ":") {
				pathParams[routePart[1:]] = pathPart
			} else {
				found = false
				break
			}
		}

		if found {
			return routable, pathParams, route
		}
	}

	// Not found
	return nil, params, ""
}

func getParts(path string) []string {
	return strings.Split(path, "/")
}
