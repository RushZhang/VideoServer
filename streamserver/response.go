package main

import (
	"net/http"
	"io"
)

func sendErrorResonse(w http.ResponseWriter, sc int, errMsg string) {
	w.WriteHeader(sc)  //返回码是写header

	io.WriteString(w, errMsg)
}