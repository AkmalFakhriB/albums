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
	http.HandleFunc("/albums/updateprice/{id}", API.UpdateAlbumPrice)
	http.HandleFunc("/albums/delete/{id}", API.DeleteAlbum)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
