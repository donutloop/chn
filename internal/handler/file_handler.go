package handler

import "net/http"

// fileHandler serves a file like the favicon or logo
func File(file string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch file {
		case "favicon":
			http.ServeFile(w, r, "../../static/img/favicon.ico")
		case "logo":
			http.ServeFile(w, r, "../../static/img/logo.gif")
		case "index":
			http.ServeFile(w, r, "../../static/html/index.html")
		case "bundle.js":
			http.ServeFile(w, r, "../../static/js/bundle.js")
		case "main.css":
			http.ServeFile(w, r, "../../static/css/main.css")
		default:
			w.WriteHeader(404)
			w.Write([]byte("file not found"))
		}
	}
}
