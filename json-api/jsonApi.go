package jsonapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func fromJSON[T any](body io.Reader, target T) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	json.Unmarshal(buf.Bytes(), &target)

}

func returnJSON[T any](w http.ResponseWriter, withData func() (T, error)) {
	setJSONHeader(w)

	data, serverErr := withData()

	if serverErr != nil {
		w.WriteHeader(500)
		serverErrJson, err := json.Marshal(&serverErr)
		if err != nil {
			log.Println(err)
			return
		}
		w.Write(serverErrJson)
		return
	}
	dataJson, err := json.Marshal(&data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(dataJson)
}

func returnErr(w http.ResponseWriter, err error, code int) {
	returnJSON(w, func() (interface{}, error) {
		errorMessage := struct {
			Err string
		}{
			Err: err.Error(),
		}
		w.WriteHeader(code)
		return errorMessage, nil
	})

}
