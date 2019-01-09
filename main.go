package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
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

func indexHandler(w http.ResponseWriter, r *http.Request) error {
	fd, err := os.Open("index.html")
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = io.Copy(w, fd)
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
