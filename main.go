package main

import (
	_ "embed"
	"errors"
	"log"
	"net/http"
	"time"
)

func latexHandler(w http.ResponseWriter, r *http.Request) error {
	defer func(t time.Time) {
		log.Println("execution time", time.Since(t))
	}(time.Now())

	mr, err := r.MultipartReader()
	if err != nil {
		return err
	}
	return printMultipart(w, mr)
}

//go:embed index.html
var index []byte

func indexHandler(w http.ResponseWriter, r *http.Request) error {
	_, err := w.Write(index)
	return err
}

func errHandler(w http.ResponseWriter, err error, code int) {
	if err == nil {
		err = errors.New(http.StatusText(code))
	}
	http.Error(w, err.Error(), code)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := indexHandler(w, r); err != nil {
			errHandler(w, err, http.StatusInternalServerError)
		}
	case http.MethodPost:
		if err := latexHandler(w, r); err != nil {
			errHandler(w, err, http.StatusInternalServerError)
		}
	default:
		errHandler(w, nil, http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
