package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	e := New()
	srv := httptest.NewServer(e)

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatalf("error http GET %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("statusCode expected %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
