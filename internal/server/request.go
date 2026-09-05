package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func DecodeRequestBody(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode JSON body: %w", err)
	}
	return nil
}
