package db

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
)

type Albums struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func ConnectDB() (*sql.DB, error) {
	connStr := "user=Akmal password=4kum4RUp0stgr3s dbname=recordings sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func AddAlbum(album Albums) (int64, error) {
	var newId int64
	db, err := ConnectDB()
	if err != nil {
		fmt.Printf("Error when trying to connect to DB: %s", err)
		return 0, err
	}

	err = db.QueryRow("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id", album.Title, album.Artist, album.Price).Scan(&newId)
	if err != nil {
		return 0, err
	}

	return newId, nil
}

func AlbumsByMinimumPrice(price float32) ([]Albums, error) {
	var albums []Albums
	db, err := ConnectDB()
	if err != nil {
		fmt.Printf("Error when trying to connect to DB: %s", err)
		return albums, err
	}

	rows, err := db.Query(`SELECT id, title, artist, price FROM album WHERE price >= $1`, price)
	defer rows.Close()

	for rows.Next() {
		alb := Albums{}

		s := reflect.ValueOf(&alb).Elem()
		numCols := s.NumField()
		columns := make([]interface{}, numCols)
		for i := 0; i < numCols; i++ {
			field := s.Field(i)
			columns[i] = field.Addr().Interface()
		}

		err := rows.Scan(columns...)
		if err != nil {
			return albums, err
		}
		albums = append(albums, alb)
	}

	if err = rows.Err(); err != nil {
		return albums, err
	}

	return albums, nil
}

func AllAlbums() ([]Albums, error) {
	var albums []Albums

	db, err := ConnectDB()
	if err != nil {
		fmt.Printf("Error when trying to connect to DB: %s", err)
		return albums, err
	}

	rows, err := db.Query("SELECT id, title, artist, price FROM album")
	if err != nil {
		return albums, err
	}
	defer rows.Close()

	for rows.Next() {
		alb := Albums{}

		s := reflect.ValueOf(&alb).Elem()
		numCols := s.NumField()
		columns := make([]interface{}, numCols)
		for i := 0; i < numCols; i++ {
			field := s.Field(i)
			columns[i] = field.Addr().Interface()
		}

		err := rows.Scan(columns...)
		if err != nil {
			return albums, err
		}
		albums = append(albums, alb)
	}

	if err = rows.Err(); err != nil {
		return albums, err
	}

	return albums, nil
}

func AlbumsById(id int64) (Albums, error) {
	var album Albums
	db, err := ConnectDB()
	if err != nil {
		fmt.Printf("Error when trying to connect to DB: %s", err)
		return album, err
	}

	row := db.QueryRow(`SELECT id, title, artist, price FROM album WHERE id = $1`, id)

	err = row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
	if err != nil {
		return album, err
	}

	return album, nil
}

func ChangeAlbumPrice(id int64, newPrice float32) (int64, error) {
	db, err := ConnectDB()
	if err != nil {
		fmt.Printf("Error when trying to connect to DB: %s", err)
		return 0, err
	}

	_, err = db.Exec("UPDATE album SET price = $1 WHERE id = $2", newPrice, id)
	if err != nil {
		fmt.Printf("Error when trying to update data to DB: %s", err)
		return 0, err
	}

	return id, nil
}

func DeleteAlbum(id int64) (int64, error) {
	db, err := ConnectDB()
	if err != nil {
		fmt.Printf("Error when trying to connect to DB: %s", err)
		return 0, err
	}

	res, err := db.Exec("DELETE FROM album WHERE id = $1", id)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}
