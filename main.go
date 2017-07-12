package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func recvFiles(dir string, r *multipart.Reader) error {
	part, err := r.NextPart()
	if err != nil {
		return err
	}
	defer part.Close()
	fname := part.FileName()
	if fname == "" {
		return nil
	}
	log.Println("recv", fname)
	fd, err := os.Create(filepath.Join(dir, filepath.Base(fname)))
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = io.Copy(fd, part)
	return err
}

func texHandler(w http.ResponseWriter, r *http.Request) error {
	defer func(t time.Time) {
		log.Println("execution time", time.Since(t))
	}(time.Now())
	dir, err := ioutil.TempDir("", "gotex")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	mr, err := r.MultipartReader()
	if err != nil {
		return err
	}
	for {
		err := recvFiles(dir, mr)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	out := new(bytes.Buffer)
	cmd := exec.Command("latexmk", "-pdf", "-interaction=nonstopmode", "-jobname=out")
	cmd.Dir = dir
	cmd.Stdout = out
	cmd.Stderr = out

	if err := cmd.Run(); err != nil {
		return errors.New(out.String() + err.Error())
	}

	return sendFile(filepath.Join(dir, "out.pdf"), w, r)
}

func sendFile(fname string, w http.ResponseWriter, r *http.Request) error {
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
		if err := sendFile("index.html", w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		if err := texHandler(w, r); err != nil {
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
