package main

import (
	API "example/postgrest1/api"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", API.RedirectHandler)
	http.HandleFunc("/albums", API.GetAllAlbums)
	http.HandleFunc("/albums/{id}", API.GetAlbumById)
	http.HandleFunc("/albums/newalbum", API.CreateNewAlbum)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
