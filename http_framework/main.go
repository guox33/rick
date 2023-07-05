package main

import (
	"log"
	"net/http"
)

type fooHandler struct {
}

func (f *fooHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	_, _ = writer.Write([]byte("you know"))
	writer.WriteHeader(200)
}

func main() {
	http.Handle("/foo", &fooHandler{})

	http.HandleFunc("/bar", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("hello world"))
		writer.WriteHeader(200)
	})

	http.NotFoundHandler()

	log.Fatal(http.ListenAndServe(":8000", nil))
}
