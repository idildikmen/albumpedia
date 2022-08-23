package main

// The sql go library is needed to interact with the database
import (
	"database/sql"
	"log"
)

// Our store will have two methods, to add a new bird,
// and to get all existing birds
// Each method returns an error, in case something goes wrong
type Store interface {
	CreateAlbum(album *Album) error
	GetAlbums() ([]*Album, error)
}

// The `dbStore` struct will implement the `Store` interface
// It also takes the sql DB connection object, which represents
// the database connection.
type dbStore struct {
	db *sql.DB
}

func (store *dbStore) CreateAlbum(album *Album) error {
	// 'Bird' is a simple struct which has "species" and "description" attributes
	// THe first underscore means that we don't care about what's returned from
	// this insert query. We just want to know if it was inserted correctly,
	// and the error will be populated if it wasn't
	_, err := store.db.Exec("INSERT INTO albums(title, artist, price) VALUES ($1,$2,$3)", album.Title, album.Artist, album.Price)
	return err
}
func (store *dbStore) CreateTestAlbum(album *Album) error {
	// 'Bird' is a simple struct which has "species" and "description" attributes
	// THe first underscore means that we don't care about what's returned from
	// this insert query. We just want to know if it was inserted correctly,
	// and the error will be populated if it wasn't
	_, err := store.db.Exec("INSERT INTO testalbum(title, artist, price) VALUES ($1,$2,$3)", album.Title, album.Artist, album.Price)
	return err
}

func (store *dbStore) GetAlbums() ([]*Album, error) {
	// Query the database for all birds, and return the result to the
	// `rows` object
	rows, err := store.db.Query("SELECT title, artist, price from albums")
	// We return incase of an error, and defer the closing of the row structure
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create the data structure that is returned from the function.
	// By default, this will be an empty array of birds
	albums := []*Album{}
	for rows.Next() {
		// For each row returned by the table, create a pointer to a bird,
		album := &Album{}
		// Populate the `Species` and `Description` attributes of the bird,
		// and return incase of an error
		if err := rows.Scan(&album.Title, &album.Artist, &album.Price); err != nil {
			return nil, err
		}
		// Finally, append the result to the returned array, and repeat for
		// the next row
		albums = append(albums, album)
	}
	return albums, nil
}

func (store *dbStore) GetTestAlbums() ([]*Album, error) {
	// Query the database for all birds, and return the result to the
	// `rows` object
	rows, err := store.db.Query("SELECT title, artist, price from testalbum")
	// We return incase of an error, and defer the closing of the row structure
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create the data structure that is returned from the function.
	// By default, this will be an empty array of birds
	albums := []*Album{}
	for rows.Next() {
		// For each row returned by the table, create a pointer to a bird,
		album := &Album{}
		// Populate the `Species` and `Description` attributes of the bird,
		// and return incase of an error
		if err := rows.Scan(&album.Title, &album.Artist, &album.Price); err != nil {
			return nil, err
		}
		// Finally, append the result to the returned array, and repeat for
		// the next row
		albums = append(albums, album)
	}
	return albums, nil
}

// The store variable is a package level variable that will be available for
// use throughout our application code
var store Store

/*
We will need to call the InitStore method to initialize the store. This will
typically be done at the beginning of our application (in this case, when the server starts up)
This can also be used to set up the store as a mock, which we will be observing
later on
*/
func InitStore(s Store) {
	store = s
}

func createTestTable(db *sql.DB) {

	dropTableSQL := `DROP TABLE IF EXISTS testalbum`
	st, err := db.Prepare(dropTableSQL)
	st.Exec()
	createAlbumsTableSQL := `CREATE TABLE testalbum (
		"idAlbum" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"title" TEXT,
		"artist" TEXT,
		"price" TEXT,
		"class" TEXT,
		"genre" TEXT
	  );` // SQL Statement for Create Table

	statement, err := db.Prepare(createAlbumsTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}

	statement.Exec() // Execute SQL Statements

}
