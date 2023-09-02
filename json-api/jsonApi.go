package jsonapi

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mailinglist/mdb"
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
func CreateEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			return
		}
		entry := mdb.EmailEntry{}
		fromJSON(req.Body, &entry)
		if err := mdb.CreateEmail(db, entry.Email); err != nil {
			returnErr(w, err, 400)
			return
		}
		returnJSON(w, func() (interface{}, error) {
			log.Printf("JSON CreateEmail : %v\n", entry.Email)
			return mdb.GetEmail(db, entry.Email)
		})

	})
}
func GetEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			return
		}
		entry := mdb.EmailEntry{}
		fromJSON(req.Body, &entry)

		returnJSON(w, func() (interface{}, error) {
			log.Printf("JSON GetEmail : %v\n", entry.Email)
			return mdb.GetEmail(db, entry.Email)
		})

	})
}

func GetEmailBatch(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			return
		}
		queryOptions := mdb.BatchEmailQueryParams{}
		fromJSON(req.Body, &queryOptions)

		if queryOptions.Count <= 0 || queryOptions.Page <= 0 {
			returnErr(w, errors.New("page and count fields are required and must be > 0"), 400)
			return
		}

		returnJSON(w, func() (interface{}, error) {
			log.Printf("JSON GetEmailBatch: %v\n", queryOptions)
			return mdb.GetEmailBatch(db, queryOptions)
		})
	})
}

func UpdateEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "PATCH" {
			return
		}
		entry := mdb.EmailEntry{}
		fromJSON(req.Body, &entry)
		if err := mdb.UpdateEmail(db, entry); err != nil {
			returnErr(w, err, 400)
			return
		}
		returnJSON(w, func() (interface{}, error) {
			log.Printf("JSON UpdateEmail : %v\n", entry.Email)
			return mdb.GetEmail(db, entry.Email)
		})

	})
}

func DeleteEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "DELETE" {
			return
		}
		entry := mdb.EmailEntry{}
		fromJSON(req.Body, &entry)
		if err := mdb.DeleteEmail(db, entry.Email); err != nil {
			returnErr(w, err, 400)
			return
		}
		returnJSON(w, func() (interface{}, error) {
			log.Printf("JSON DeleteEmail : %v\n", entry.Email)
			return mdb.GetEmail(db, entry.Email)
		})

	})
}

func Serve(db *sql.DB, bind string) {
	http.Handle("/email/create", CreateEmail(db))
	http.Handle("/email/get", GetEmail(db))
	http.Handle("/email/get_batch", GetEmailBatch(db))
	http.Handle("/email/update", UpdateEmail(db))
	http.Handle("/email/delete", DeleteEmail(db))
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		log.Fatalf("json server error : %v\n", err)
		return
	}
}
