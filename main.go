package main

import (
	"fmt"
	"net/http"
	"log"
	"strings"
	"encoding/gob"
)


type Handler struct {
	fileServer http.Handler
}


func main() {
	gob.Register(sesKey(0))

	// Resource files routs (js, css)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./scripts"))))

	fmt.Println("WebServer started OK")
	fmt.Println("Try http://localhost:8085")
	//fmt.Println("or https://localhost:443")

	http.ListenAndServe(":8085", &Handler{
		fileServer: http.FileServer(http.Dir("www")),
	})
}


// Роутер
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v", r.Method, r.URL.Path)
	// need to serve shell via websocket?
	if strings.Trim(r.URL.Path, "/") == "shell" {
		onShell(w, r)
		return
	}
	
	if strings.Trim(r.URL.Path, "/") == "" {
		chrony_page(w, r)
		isStartChrony()
		return
	}
	
	if strings.Trim(r.URL.Path, "/") == "poll" {
		PollResponse(w, r)
		return
	}
	
	if strings.Trim(r.URL.Path, "/") == "is_act" {
		PollResponse2(w, r)
		return
	}
	
	if strings.Trim(r.URL.Path, "/") == "login" {
		login(w, r)
		return
	}
	
	if strings.Trim(r.URL.Path, "/") == "logout" {
		logout(w, r)
		return
	}
	
	
		
	// serve static assets from 'static' dir:
	h.fileServer.ServeHTTP(w, r)
}


