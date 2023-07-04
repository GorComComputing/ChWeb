package main

import (
	"fmt"
	"net/http"
	"log"
	"strings"
)


type Handler struct {
	fileServer http.Handler
}


func main() {
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
	
	if strings.Trim(r.URL.Path, "/") == "sh" {
		shell(w, r)
		isStartChrony()
		return
	}
	// Запуск Chrony
	if strings.Trim(r.URL.Path, "/") == "start" {
		start(w, r)
		return
	}
	// Остановка Chrony
	if strings.Trim(r.URL.Path, "/") == "stop" {
		stop(w, r)
		return
	}
	// Перезапуск Chrony
	if strings.Trim(r.URL.Path, "/") == "restart" {
		restart(w, r)
		return
	}
	
	if strings.Trim(r.URL.Path, "/") == "save" {
		save(w, r)
		return
	}
	
	if strings.Trim(r.URL.Path, "/") == "activity" {
		activity(w, r)
		return
	}
	
	if strings.Trim(r.URL.Path, "/") == "tracking" {
		tracking(w, r)
		return
	}
	
	if strings.Trim(r.URL.Path, "/") == "sources" {
		sources(w, r)
		return
	}
	
	if strings.Trim(r.URL.Path, "/") == "sourcestats" {
		sourcestats(w, r)
		return
	}
	
	if strings.Trim(r.URL.Path, "/") == "clients" {
		clients(w, r)
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
	
	if strings.Trim(r.URL.Path, "/") == "config" {
		config(w, r)
		return
	}
	
	if strings.Trim(r.URL.Path, "/") == "saveconfig" {
		saveconfig(w, r)
		return
	}
	
	if strings.Trim(r.URL.Path, "/") == "restore" {
		restore(w, r)
		return
	}
	
	
	// serve static assets from 'static' dir:
	h.fileServer.ServeHTTP(w, r)
}


