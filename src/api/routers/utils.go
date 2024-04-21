package routers

import (
	"encoding/json"
	"net/http"
)

// respondWithError sets the error response to the response writer
func respondWithError(res http.ResponseWriter, code int, message string) {
	respondWithJSON(res, code, map[string]string{"error": message})
}

// respondWithJSON sets the response payload in the response writer
func respondWithJSON(res http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	res.Write(response)
}
