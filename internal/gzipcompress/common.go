package gzipcompress

import (
	"net/http"
	"slices"
	"strings"
)

const FormatName = "gzip"

func IsDecompressValid(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), FormatName)
}

func IsCompressValid(rw http.ResponseWriter, r *http.Request) bool {
	isAccept := r.Header.Get("Content-Encoding") == FormatName

	validRespContentType := []string{"application/json", "text/html"}
	isValidResponseType := slices.Contains(validRespContentType, rw.Header().Get("Content-Type"))

	return isAccept && isValidResponseType
}
