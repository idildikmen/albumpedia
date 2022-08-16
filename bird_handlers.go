package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Album struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Price  string `json:"price"`
}

var albums []Album

func getAlbumHandler(w http.ResponseWriter, r *http.Request) {
	//Convert the "birds" variable to json
	albumListBytes, err := json.Marshal(albums)

	// If there is an error, print it to the console, and return a server
	// error response to the user
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// If all goes well, write the JSON list of birds to the response
	w.Write(albumListBytes)
}

func createAlbumHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new instance of Bird
	album := Album{}

	// We send all our data as HTML form data
	// the `ParseForm` method of the request, parses the
	// form values
	err := r.ParseForm()

	// In case of any error, we respond with an error to the user
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the information about the bird from the form info
	album.Title = r.Form.Get("title")
	album.Artist = r.Form.Get("artist")
	album.Price = r.Form.Get("price")

	// Append our existing list of birds with a new entry
	albums = append(albums, album)

	//Finally, we redirect the user to the original HTMl page
	// (located at `/assets/`), using the http libraries `Redirect` method
	http.Redirect(w, r, "/assets/", http.StatusFound)
}
