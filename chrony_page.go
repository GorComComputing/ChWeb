package main

import (
	"html/template"
	"net/http"	
 	"io"
)


// Page fields
type ShellPage struct {
	Message  	string
	UserName 	string
	BackLink 	string
	Conf	 	Config
	IsStarted	bool
	File		string //template.HTML 
}

var bks ShellPage


// создание каналов для long poll
var messages chan string = make(chan string, 100)
var messages2 chan string = make(chan string, 100)


// получение longpoll и установка канала для ответа
func PollResponse(w http.ResponseWriter, req *http.Request) {
    	io.WriteString(w, <-messages)
}


// получение longpoll и установка канала для ответа
func PollResponse2(w http.ResponseWriter, req *http.Request) {
    	io.WriteString(w, <-messages2)
}


// Chrony page handler
func chrony_page(w http.ResponseWriter, r *http.Request) {
	
	// Set page fields
	bks.Message = "Message"
	bks.UserName = "User"//check_cookies(w, r)
	bks.BackLink = "sh"
	
	bks.Conf, bks.File = scan()
	//File = strings.Replace(File, "\n", "<br/>", -1)
	//bks.File = template.HTML(File)
	//fmt.Println(bks.File )
	
	bks.IsStarted = isStartChrony()

	// Response template
	tmpl, _ := template.ParseFiles("templates/shell.html")
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, bks)
}


