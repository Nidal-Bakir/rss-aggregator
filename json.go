package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const ContentTypeJSON = "application/json"

func responseWithError(w http.ResponseWriter, code int, msg string, err error) {
	if code > 499 {
		fmt.Printf("Internal server Error. Code:%d , Message: %s, Error: %v", code, msg, err)
	}

	errorJson := struct {
		ErrorMsg string `json:"error"`
	}{ErrorMsg: msg}

	responseWithJson(w, code, errorJson)
}

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error while marshaling the payload:", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(code)
	w.Write(data)
}
