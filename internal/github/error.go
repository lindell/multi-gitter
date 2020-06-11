package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func responseToError(r *http.Response, msg string) error {
	var resp errorResponse
	_ = json.NewDecoder(r.Body).Decode(&resp)
	if resp.Message != "" {
		return fmt.Errorf("%s, got status code %d with message: %s", msg, r.StatusCode, resp.Message)
	}
	return fmt.Errorf("%s, got status code %d", msg, r.StatusCode)
}
