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

func TestGetAlbumsHandler(t *testing.T) {

	albums = []Album{
		{"New", "Imagine Dragons", "22.99"},
	}

	req, err := http.NewRequest("GET", "", nil)

	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(getAlbumHandler)

	hf.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := Album{"New", "Imagine Dragons", "22.99"}
	b := []Album{}
	err = json.NewDecoder(recorder.Body).Decode(&b)

	if err != nil {
		t.Fatal(err)
	}

	actual := b[0]

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}
func TestCreateAlbumsHandler(t *testing.T) {

	albums = []Album{
		{"New", "Imagine Dragons", "22.99"},
	}

	form := newCreateAlbumForm()
	req, err := http.NewRequest("POST", "", bytes.NewBufferString(form.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(createAlbumHandler)

	hf.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := Album{"Halo", "Beyonce", "33.99"}

	if err != nil {
		t.Fatal(err)
	}

	actual := albums[1]

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func newCreateAlbumForm() *url.Values {
	form := url.Values{}
	form.Set("title", "Halo")
	form.Set("artist", "Beyonce")
	form.Set("price", "33.99")
	return &form
}
