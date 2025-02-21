package routes

import (
	"io/fs"
	"log"
	"net/http"
	"ppl-calculations/ports/http/middleware"
)

func RegisterStaticRoutes(mux *http.ServeMux, assets fs.FS) {
	cssFs, err := fs.Sub(assets, "assets/css")
	if err != nil {
		log.Fatalf("Error parsing css: %v", err)
	}
	mux.Handle("/css/", middleware.CacheHeaders(http.StripPrefix("/css/", http.FileServer(http.FS(cssFs)))))

	imagesFs, err := fs.Sub(assets, "assets/images")
	if err != nil {
		log.Fatalf("Error parsing images: %v", err)
	}
	mux.Handle("/images/", middleware.CacheHeaders(http.StripPrefix("/images/", http.FileServer(http.FS(imagesFs)))))

	jsFs, err := fs.Sub(assets, "assets/js")
	if err != nil {
		log.Fatalf("Error parsing js: %v", err)
	}
	mux.Handle("/js/", http.StripPrefix("/js/", middleware.CacheHeaders(http.FileServer(http.FS(jsFs)))))

	fontFs, err := fs.Sub(assets, "assets/fonts")
	if err != nil {
		log.Fatalf("Error parsing fonts: %v", err)
	}
	mux.Handle("/fonts/", http.StripPrefix("/fonts/", middleware.CacheHeaders(http.FileServer(http.FS(fontFs)))))
}
