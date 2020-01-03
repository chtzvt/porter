package api

import (
	"net/http"
	"sort"
	"strings"
)

// genKeyMap converts an array of config.APIKey into a map of tokens ("Secrets") to
// the whitelisted methods they are allowed to access.
// This map generated once at startup and used to simplify permissions lookups.
func (s *Server) genKeyMap() {
	for _, apiKey := range s.Config.Keys {
		s.KeyMap[apiKey.Secret] = apiKey.MethodWhitelist
	}
}

// checkAuth enforces permissions on incoming HTTP requests
// For API tokens provided in an incoming HTTP request, checkAuth ensures that the
// resource being accessed is whitelisted.
func (s *Server) checkAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Authorization is disabled because no API keys are configured
		if len(s.KeyMap) == 0 {
			h(w, r)
			return
		}

		apiKey := r.Header.Get("Authorization")

		// Empty API key provided in request
		if apiKey == "" {
			http.Error(w, NewStatusString("error", "missing API key"), http.StatusForbidden)
			return
		}

		token := strings.Replace(apiKey, "Bearer ", "", 1)
		allowedMethods, validToken := s.KeyMap[token]

		// API key is invalid or exists with no whitelisted methods defined
		if !validToken || len(allowedMethods) == 0 {
			http.Error(w, NewStatusString("error", "unauthorized"), http.StatusForbidden)
			return
		}

		// Provided API Key is allowed access to any method
		if allowedMethods[0] == "*" {
			h(w, r)
			return
		}

		method := getMethodName(r.URL.Path)

		// Requested method not in API Key whitelist
		if i := sort.SearchStrings(allowedMethods, method); i >= len(allowedMethods) || allowedMethods[i] != method {
			http.Error(w, NewStatusString("error", "unauthorized"), http.StatusForbidden)
			return
		}

		h(w, r)
	}
}

func getMethodName(path string) string {
	basePath := strings.Replace(path, "/api/v1/", "", 1)
	if trailingSlash := strings.LastIndex(basePath, "/"); trailingSlash != -1 {
		return basePath[:trailingSlash]
	}

	return basePath
}
