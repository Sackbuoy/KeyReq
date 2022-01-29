package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"time"
	"math/rand"
)

type Key struct {
	Value string `json:"value"`
	Issued int64 `json:"issued"`
	Expiry int64 `json:"expiry"`
}

type ValidationResponse struct {
	IsValid bool `json:"isValid"`
}

const PORT = ":8000"
var validKeys []Key

func generateKey() string {
	charset := "abcdefghijklmnopqrstuvwxyz" +
	  "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))


	b := make([]byte, 128)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

const getKeyPath string = "/getKey"
func getKey(writer http.ResponseWriter, req *http.Request) {
	// get Unix timestamp for the issued time, key is valid for 1 hour, and generate the key
	issued := time.Now().Unix()
	expiry := issued + 3600
	keyValue := generateKey()

	key := Key{Value: keyValue, Issued: issued, Expiry: expiry}
	jsonKey, err := json.Marshal(key)
	if err != nil {
		panic(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonKey)

	validKeys = append(validKeys, key)
}

const validateKeyPath string = "/validate/{key}"
func validateKey(writer http.ResponseWriter, req *http.Request) { 
	vars := mux.Vars(req)
	key := vars["key"]

	response := ValidationResponse{IsValid: false}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	for _, validKey := range validKeys {
		if key == validKey.Value && validKey.Expiry > time.Now().Unix() {
			response.IsValid = true
		}
	}	

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	writer.Write(jsonResponse)
}

func main() {
	fmt.Println("Server Started")
	router := mux.NewRouter()
	router.HandleFunc(getKeyPath, getKey).Methods("GET")
	router.HandleFunc(validateKeyPath, validateKey).Methods("GET")

	log.Fatal(http.ListenAndServe(PORT, router))
}
