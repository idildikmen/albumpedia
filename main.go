package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/mux"
)

// The new router function creates the router and
// returns it to us. We can now use this function
// to instantiate and test the router outside of the main function
func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hello", handler).Methods("GET")
	// Declare the static file directory and point it to the
	// directory we just made
	staticFileDirectory := http.Dir("./assets/")
	// Declare the handler, that routes requests to their respective filename.
	// The fileserver is wrapped in the `stripPrefix` method, because we want to
	// remove the "/assets/" prefix when looking for files.
	// For example, if we type "/assets/index.html" in our browser, the file server
	// will look for only "index.html" inside the directory declared above.
	// If we did not strip the prefix, the file server would look for
	// "./assets/assets/index.html", and yield an error
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	// The "PathPrefix" method acts as a matcher, and matches all routes starting
	// with "/assets/", instead of the absolute route itself
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")
	// These lines are added inside the newRouter() function before returning r
	r.HandleFunc("/album", getAlbumHandler).Methods("GET")
	r.HandleFunc("/album", createAlbumHandler).Methods("POST")
	return r
}

// Users struct which contains
// an array of users
type Albums struct {
	Albums []album `json:"albums"`
}

// User struct which contains a name
// a type and a list of social links
//album represents data about a record album
type album struct {
	Class  string `json:"_class"`
	Artist string `json:"artist"`
	Title  string `json:"title"`
	Year   string `json:"releaseYear"`
	Genre  string `json:"genre"`
}

func main() {
	// The router is now formed by calling the `newRouter` constructor function
	// that we defined above. The rest of the code stays the same

	jsonFile, err := os.Open("albums.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened albums.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var albums Albums

	json.Unmarshal(byteValue, &albums)
	fmt.Println("First album title is:", albums.Albums[0].Title)

	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database-alb.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database-alb.db")
	defer sqliteDatabase.Close() // Defer Closing the database
	createTable(sqliteDatabase)

	for i := 0; i < len(albums.Albums); i++ {
		album := albums.Albums[i]
		_, err1 := sqliteDatabase.Exec("INSERT INTO albums(title, artist, class, genre, year) VALUES ($1,$2,$3,$4,$5)", album.Title, album.Artist,
			album.Class, album.Genre, album.Year)
		_ = err1
	}

	InitStore(&dbStore{db: sqliteDatabase})

	r := newRouter()
	log.Println("comes here")
	fmt.Println("Serving on port 8080")
	http.ListenAndServe(":8080", r)
}

func createTable(db *sql.DB) {
	createAlbumsTableSQL := `CREATE TABLE albums (
		"idAlbum" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"title" TEXT,
		"artist" TEXT,
		"price" TEXT,
		"class" TEXT,
		"genre" TEXT,
		"year" TEXT
	  );` // SQL Statement for Create Table

	log.Println("Create album table...")
	statement, err := db.Prepare(createAlbumsTableSQL) // Prepare SQL Statement
	log.Println("now")
	if err != nil {
		log.Fatal(err.Error())
	}

	statement.Exec() // Execute SQL Statements
	log.Println("album table created")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hey it's Idil!")

}
