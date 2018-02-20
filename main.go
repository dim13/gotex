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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := indexHandler(w, r); err != nil {
			code := http.StatusInternalServerError
			http.Error(w, err.Error(), code)
		}
	case http.MethodPost:
		if err := latexHandler(w, r); err != nil {
			code := http.StatusInternalServerError
			http.Error(w, err.Error(), code)
		}
	default:
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
