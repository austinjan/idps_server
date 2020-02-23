package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	mongodb "github.com/austinjan/idps_server/servers/mongo"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// Get JSON body and transform to bson
func parseBody2Bson(r *http.Request) (map[string]interface{}, error) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		return nil, err
	}

	var bodyBson map[string]interface{}

	if err := json.Unmarshal(b, &bodyBson); err != nil {
		return nil, err
	}

	return bodyBson, err
}

// api rul format : /apis/[version]/[collection]?[query]
// return body: array of json for data
func readData(w http.ResponseWriter, r *http.Request) {

	//_, collection, err := bodyParser(r)

	vars := mux.Vars(r)
	collection := vars["collection"]
	fmt.Println("/v1/read/", collection)
	query, err := createBsonQuery(r)
	if err != nil {
		handleResponse(w, []byte(err.Error()), err)
		return
	}

	from, fromerr := query["from"]
	to, toerr := query["to"]

	if fromerr != false && toerr != false {
		_from,_ := strconv.Atoi(from.(string))
		_to,_ := strconv.Atoi(to.(string))

		query = bson.M{"startTime": bson.M{"$gte": _from, "$lt": _to}}
		//query = bson.M{"startTime": bson.M{"$gte": 1582041600000}}
	}

	db := mongodb.GetDB()
	jsonRet, err := db.Read(collection, query)
	if err != nil {
		handleResponse(w, []byte(err.Error()), err)
		return
	}
	handleResponse(w, jsonRet, err)
	// handleResponse(w, []byte("hello"), nil)
}

func testReport(w http.ResponseWriter, r *http.Request) {

	bodyBson, err := parseBody2Bson(r)

	if err != nil {
		handleResponse(w, nil, err)
		return
	}
	//fmt.Println("POST apis ", bodyBson)

	db := mongodb.GetDB()
	err = db.Insert("tagtesting", bodyBson)
	if err != nil {
		handleResponse(w, []byte(err.Error()), err)
		return
	}

	jsonStr, err := json.Marshal(bodyBson)
	if err != nil {
		handleResponse(w, []byte(err.Error()), err)
		return

	}

	handleResponse(w, jsonStr, err)

	return
}

func initCrudRouter(r *mux.Router) {
	r.HandleFunc("/api/{version:v[0-9]+}/read/{collection}", readData).Methods("GET")
	r.HandleFunc("/api/testreport", testReport).Methods("POST")
}

// InitRouter  Init rootRouter
func InitRouter(r *mux.Router) {
	initCrudRouter(r)
}
