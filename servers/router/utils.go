package router

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

// ErrorResponse error response
type ErrorResponse struct {
	Error string `json:"error"`
}

func createBsonQuery(r *http.Request) (map[string]interface{}, error) {
	querys, err := url.ParseQuery(r.URL.RawQuery)

	query := make(map[string]interface{})
	for k, v := range querys {
		// Right now I just recieve 1 parameters, will implement to support multiple parameters
		query[k] = v[0]
	}
	return query, err
}

func handleResponse(w http.ResponseWriter, body []byte, err error) {

	// fail
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		res := ErrorResponse{string(body)}
		json, err := json.Marshal(res)
		w.Write(json)

		log.Println("CRUD error- ", err)
		return
	}
	// success
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if body != nil {
		w.Write(body)
	}
	return
}
