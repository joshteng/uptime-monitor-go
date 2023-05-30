package utils

import "net/http"
import "encoding/json"

func ReturnJsonResponse(w http.ResponseWriter, httpCode int, response map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(response)
}
