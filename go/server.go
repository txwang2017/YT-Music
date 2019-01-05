package main

import (
	"fmt"
	"net/http"
)

type UrlHander func(w http.ResponseWriter, r *http.Request)

type Server struct {
	port int
}

func (server *Server) addUrl(url string, handler UrlHander) {
	http.HandleFunc(url, handler)
}

func (server *Server) start() {
	http.ListenAndServe(fmt.Sprintf(":%d", server.port), nil)
}

func StartServer() {
	server := Server{port: 8000}
	server.addUrl("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
	server.start()
}
