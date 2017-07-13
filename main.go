package main

import (
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
	return Print(w, mr)
}

func indexHandler(fname string, w http.ResponseWriter, r *http.Request) error {
	fd, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = io.Copy(w, fd)
	return err
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := indexHandler("index.html", w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		if err := latexHandler(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(":8080", nil)
}
