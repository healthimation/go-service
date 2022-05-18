package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/healthimation/go-glitch/glitch"
)

// ReturnProblem will return a json http problem response
func ReturnProblem(w http.ResponseWriter, detail, code string, status int) (int, []byte) {
	prob := glitch.HTTPProblem{
		Title:  http.StatusText(status),
		Detail: detail,
		Code:   code,
		Status: status,
	}
	by, _ := json.Marshal(prob)
	if w != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	return status, by
}

// WriteProblem will write a json http problem response
func WriteProblem(w http.ResponseWriter, detail, code string, status int) error {
	log.Println(code, detail, status)
	prob := glitch.HTTPProblem{
		Title:  http.StatusText(status),
		Detail: detail,
		Code:   code,
		Status: status,
	}
	by, err := json.Marshal(prob)
	if err != nil {
		return err
	}
	if w != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		_, err = w.Write(by)
	}
	return err
}

// WriteJSONResponse will write a json response to the htt.ResponseWriter
func WriteJSONResponse(w http.ResponseWriter, status int, data interface{}) error {
	var by []byte
	var err error
	if data != nil {
		by, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}
	if w != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		_, err = w.Write(by)
	}
	return err
}
