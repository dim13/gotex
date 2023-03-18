package main

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
)

func execute(w io.Writer, dir string) error {
	out := new(bytes.Buffer)
	cmd := exec.Command("latexmk", "-interaction=nonstopmode", "-jobname=out")
	cmd.Dir = dir
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		return errors.New(out.String() + err.Error())
	}

	// copy out
	fd, err := os.Open(filepath.Join(dir, "out.pdf"))
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = io.Copy(w, fd)
	return err
}

func savePart(r *multipart.Reader, dir string) error {
	part, err := r.NextPart()
	if err != nil {
		return err
	}
	defer part.Close()

	fname := part.FileName()
	if fname == "" {
		return errors.New("no filename specified")
	}
	fd, err := os.Create(filepath.Join(dir, filepath.Base(fname)))
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = io.Copy(fd, part)
	return err
}

func saveAllParts(r *multipart.Reader, dir string) error {
	for {
		err := savePart(r, dir)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func printMultipart(w io.Writer, r *multipart.Reader) error {
	// create temporary working directory
	dir, err := os.MkdirTemp("", "gotex")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	if err := saveAllParts(r, dir); err != nil {
		return err
	}

	return execute(w, dir)
}
