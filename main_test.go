package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSimple(t *testing.T) {
	srv := httptest.NewServer(Router())
	defer srv.Close()

	data := `{"name":"fred","email":"frederic@toggl.com"}`
	resp, err := http.Post(srv.URL + "/", "application/json", bytes.NewBufferString(data))
	log.Print(srv.URL)

	if err != nil {
		t.Fatalf("got an error sending the request: %s", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("received status code: %d", resp.StatusCode)
	}
}

func TestValidation(t *testing.T) {
	srv := httptest.NewServer(Router())
	defer srv.Close()

	data := `{"name":"fred","email":"frederic"}`
	resp, err := http.Post(srv.URL + "/", "application/json", bytes.NewBufferString(data))
	log.Print(srv.URL)

	if err != nil {
		t.Fatalf("got an error sending the request: %s", err)
	}

	if resp.StatusCode != 400 {
		t.Fatalf("received status code: %d", resp.StatusCode)
	}
}
