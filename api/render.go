package api

import (
	"encoding/json"
	"net/http"
)

// Render is the simple API helper to render JSON and other data types
type Render struct {
	R *http.Request
	W http.ResponseWriter
}

// JSON renders the provided value as a JSON encoded value into the ResponseWriter.
func (r *Render) JSON(v interface{}) {
	if err := json.NewEncoder(r.W).Encode(v); err != nil {
		http.Error(r.W, "Error encoding response", http.StatusInternalServerError)
	}
}
