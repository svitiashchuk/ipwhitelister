package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed assets
var staticFS embed.FS

// StaticHandler returns an http.Handler that serves static files from the "assets/static" directory.
func StaticHandler() http.Handler {
	httpFS, err := fs.Sub(staticFS, "assets")
	if err != nil {
		// TODO: Log the error or handle it appropriately
		log.Println(err)
	}

	return http.FileServer(http.FS(httpFS))
}
