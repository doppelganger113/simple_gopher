package openapi

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(
			path,
			http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP,
		)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(root).ServeHTTP(w, r)
	})
}
