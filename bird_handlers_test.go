package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestCreateAlbumsHandler(t *testing.T) {

	// Setup request handlers
	recorder := httptest.NewRecorder()
	create_hf := http.HandlerFunc(createAlbumHandler)
	get_hf := http.HandlerFunc(getAlbumHandler)

	// Create album (beyonce)
	form := newCreateAlbumForm()
	req, err := http.NewRequest("POST", "", bytes.NewBufferString(form.Encode()))

	// Send create album request
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	if err != nil {
		t.Fatal(err)
	}
	create_hf.ServeHTTP(recorder, req)

	// check create request answer
	if status := recorder.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusFound)
	}

	recorder = httptest.NewRecorder()

	// get albums
	req, err = http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	get_hf.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// decode get albums answers
	album_list := []Album{}
	err = json.NewDecoder(recorder.Body).Decode(&album_list)

	if err != nil {
		t.Fatal(err)
	}

	expected := Album{"Halo", "Beyonce", "33.99"}

	if album_list[0] != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", album_list[0], expected)
	}
}

func newCreateAlbumForm() *url.Values {
	form := url.Values{}
	form.Set("title", "Halo")
	form.Set("artist", "Beyonce")
	form.Set("price", "33.99")
	return &form
}
