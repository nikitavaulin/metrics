package gzipcompress

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriter(t *testing.T) {
	t.Run("empty ping response stays empty", func(t *testing.T) {
		response := httptest.NewRecorder()
		writer := NewWriter(response)
		writer.WriteHeader(http.StatusOK)
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
		if response.Body.Len() != 0 || response.Header().Get("Content-Encoding") != "" {
			t.Fatalf("unexpected gzip data or header: body=%q encoding=%q", response.Body.String(), response.Header().Get("Content-Encoding"))
		}
	})

	t.Run("json response is compressed", func(t *testing.T) {
		response := httptest.NewRecorder()
		writer := NewWriter(response)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		if _, err := writer.Write([]byte(`{"ok":true}`)); err != nil {
			t.Fatal(err)
		}
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
		if response.Header().Get("Content-Encoding") != "gzip" {
			t.Fatal("missing gzip content encoding")
		}
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(reader)
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != `{"ok":true}` {
			t.Fatalf("body = %q", body)
		}
	})
}
