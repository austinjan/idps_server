package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

// api rul format : /apis/[version]/[collection]?[query]
// return body: array of json for data
func readData(w http.ResponseWriter, r *http.Request) {
	//_, collection, err := bodyParser(r)
	// vars := mux.Vars(r)
	//collection := vars["collection"]

	//query, err := createBsonQuery(r)

	// if err != nil {
	// 	handleResponse(w, []byte(err.Error()), err)
	// 	return
	// }

	// db := mongodb.GetDB()
	// jsonRet, err := db.Read(collection, query)
	// if err != nil {
	// 	handleResponse(w, []byte(err.Error()), err)
	// 	return
	// }
	// handleResponse(w, jsonRet, err)
	handleResponse(w, []byte("hello"), nil)
}

func initCrudRouter(r *mux.Router) {

	r.HandleFunc("/{version:v[0-9]+}/read/{collection}", readData).Methods("GET")

}

// InitRouter  Init rootRouter
func InitRouter(r *mux.Router) {
	s := mux.NewRouter().PathPrefix("/apis").Subrouter().StrictSlash(true)

	initCrudRouter(s)

	// r.PathPrefix("/apis").Handler(negroni.New(
	// 	negroni.HandlerFunc(ValidateTokenMiddleWare),
	// 	negroni.Wrap(s)))

	// compatible for old api, not create

}
