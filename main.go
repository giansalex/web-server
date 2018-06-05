package main

import (
	"errors"
	"log"
	"net/http"
	"os"
)

type serveParameter struct {
	serve     string
	directory string
}

type logHandler struct {
	h http.Handler
}

func (f *logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	f.h.ServeHTTP(w, r)
}

func main() {
	parameters, err := parseArgs(os.Args)
	if err != nil {
		log.Println(err)
		return
	}

	startServer(parameters)
}

func startServer(parameter *serveParameter) {
	dir := http.FileServer(http.Dir(parameter.directory))
	mux := http.NewServeMux()
	mux.Handle("/", &logHandler{dir})

	log.Println("Listen on " + parameter.serve)
	log.Println(http.ListenAndServe(parameter.serve, mux))
}

func parseArgs(args []string) (*serveParameter, error) {
	if len(args) < 3 {
		return nil, errors.New("Missing parameters")
	}

	serve := args[1]
	path := args[2]
	if err := isValidPath(path); err != nil {
		return nil, err
	}

	result := &serveParameter{
		serve:     serve,
		directory: path,
	}

	return result, nil
}

func isValidPath(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	d, err := f.Stat()
	if err != nil {
		return err
	}

	if !d.IsDir() {
		return errors.New("Invalid Directory")
	}

	return nil
}
