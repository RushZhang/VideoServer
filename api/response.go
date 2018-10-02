package main

import (
	"net/http"
	"video_server/api/defs"
	"encoding/json"
	"io"
)

func sendErrorResponse(w http.ResponseWriter, errResp defs.ErrResponse) {
	w.WriteHeader(errResp.HttpSC)

	resStr, _ := json.Marshal(&errResp.ErrorDetail)
	io.WriteString(w, string(resStr))
}


func sendNormalResponse(w http.ResponseWriter, resp string, sc int) {
	w.WriteHeader(sc)
	io.WriteString(w, resp)
}
