package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// simple rest-api
// -gorila mux

// response JSON format
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	result, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(result)
}

// response error
func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	respondWithJSON(w, statusCode, map[string]string{"error": msg})
}

// simple struct example
type ResultStruct struct {
	Id string `json:"id"`
}

func mDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	respondWithJSON(w, http.StatusOK, ResultStruct{Id: id})
}

func mPatch(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var bodyRequest map[string]interface{}
	json.Unmarshal(reqBody, &bodyRequest)
	fmt.Println("bodyRequest: ", bodyRequest)
	// do what you need... just remember patch can modify an item not replace like the put method

	json.NewEncoder(w).Encode(bodyRequest)

}

// method Update is used for PUTs services update/replace
func Update(w http.ResponseWriter, r *http.Request) {
	var something map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&something)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, something)
	}
}

func mPut(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var bodyRequest map[string]interface{}
	json.Unmarshal(reqBody, &bodyRequest)
	fmt.Println("bodyRequest: ", bodyRequest)
	// do what you need...

	json.NewEncoder(w).Encode(bodyRequest)
}

func mGet(w http.ResponseWriter, r *http.Request) {
	// ...
	w.WriteHeader(200)
	w.Write([]byte(`{"message": "test GET method"}`))
}

// postForm parse the form request
func postForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	jsonscript := r.FormValue("script")
	fmt.Println("script: ", jsonscript)
	var script []map[string]interface{}
	err := json.Unmarshal([]byte(jsonscript), &script)
	if err != nil {
		xdata := map[string]string{"c": "1", "s": "Error", "m": "Error: JSON malformed " + fmt.Sprint(err)}
		fmt.Println("Error: ", xdata)
	} else {
		fmt.Println("script: ", script)
	}
	json.NewEncoder(w).Encode("ok")
}

func mPost(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var bodyRequest map[string]interface{}
	json.Unmarshal(reqBody, &bodyRequest)
	fmt.Println("bodyRequest: ", bodyRequest)
	// do what you need...
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(bodyRequest)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(`{"r": "pong"}`))
}

// handeResquests all end points of the api become here
func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", ping)
	router.HandleFunc("/mpost", mPost).Methods("POST")
	router.HandleFunc("/postform", postForm).Methods("POST")
	router.HandleFunc("/mget", mGet).Methods("GET")
	router.HandleFunc("/mput", mPut).Methods("PUT")
	router.HandleFunc("/update", Update).Methods("PUT")
	router.HandleFunc("/mpatch", mPatch).Methods("PATCH")
	router.HandleFunc("/mdelete", mDelete).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

//
func main() {
	fmt.Println("simple rest api ...")
	handleRequests()
}
