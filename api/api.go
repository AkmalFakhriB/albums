package api

import (
	"encoding/json"
	DB "example/postgrest1/db"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/albums", http.StatusPermanentRedirect)
}

func GetAllAlbums(w http.ResponseWriter, r *http.Request) {
	albums, err := DB.AllAlbums()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(albums)
}

func GetAlbumById(w http.ResponseWriter, r *http.Request) {
	url := r.PathValue("id")
	id, err := strconv.ParseInt(url, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	album, err := DB.AlbumsById(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(album)
}

func CreateNewAlbum(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing form data: %v", err)
		return
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 32)
	if err != nil {
		fmt.Printf("Error occured when reading price: %s", err)
		return
	}

	price32 := float32(price)

	albums := DB.Albums{
		Title:  r.FormValue("title"),
		Artist: r.FormValue("artist"),
		Price:  price32,
	}

	album, err := DB.AddAlbum(albums)
	if err != nil {
		fmt.Printf("Error while calling insert query: %s", err)
		return
	}

	w.Header().Set("Content-Type", "multipart/form-data")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(album)
}
