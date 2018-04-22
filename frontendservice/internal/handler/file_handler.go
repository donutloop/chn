package handler

import (
	"github.com/go-chi/chi"
	"net/http"
	"strings"
)

// fileHandler serves a file like the favicon or logo
func File(staticPath string ,file string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch file {
		case "index":
			http.ServeFile(w, r, staticPath + "/html/index.html")
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("file not found"))
		}
	}
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}
	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
