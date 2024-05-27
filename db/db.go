package db

import (
	"database/sql"
	"reflect"
)

type Albums struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func AddAlbum(album Albums, db *sql.DB) (int64, error) {
	var newId int64
	err := db.QueryRow("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id", album.Title, album.Artist, album.Price).Scan(&newId)
	if err != nil {
		return 0, err
	}

	return newId, nil
}

func AlbumsByMinimumPrice(price float32, db *sql.DB) ([]Albums, error) {
	var albums []Albums
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

func AllAlbums(db *sql.DB) ([]Albums, error) {
	var albums []Albums
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

func AlbumsById(id int64, db *sql.DB) (Albums, error) {
	var album Albums
	row := db.QueryRow(`SELECT id, title, artist, price FROM album WHERE id = $1`, id)

	err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
	if err != nil {
		return album, err
	}

	return album, nil
}

func ChangeAlbumPrice(id int64, db *sql.DB, newPrice float32) (int64, error) {
	_, err := db.Exec("UPDATE album SET price = $1 WHERE id = $2", newPrice, id)
	if err != nil {
		panic(err)
	}

	return id, nil
}

func DeleteAlbum(id int64, db *sql.DB) (int64, error) {
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
