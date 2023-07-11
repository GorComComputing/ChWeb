package main

import (
	"fmt"
	"html/template"
	"net/http"	
 	"io"
 	"strings"
    	"strconv"
)


// Page fields
type ShellPage struct {
	Message  	string
	UserName 	string
	BackLink 	string
	UserRole 	int
	Conf	 	Config
	IsStarted	bool
	File		string //template.HTML 
	Users  		[]UserFromDB
}






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

var bks ShellPage
	
	// Set page fields
	bks.Message = "Message"
	bks.UserName, bks.UserRole = check_cookies(w, r)
	bks.BackLink = "/"

	
	bks.Conf, bks.File = scan()
	//File = strings.Replace(File, "\n", "<br/>", -1)
	//bks.File = template.HTML(File)
	//fmt.Println(bks.File )
	
	bks.IsStarted = isStartChrony()
	
	
	openDB("http://192.168.63.60", "8086") // "http://192.168.1.136"  "http://192.168.63.60"
	req := db_request("select * from users") 

	strArr := strings.FieldsFunc(req, func(c rune) bool { return c == '\n' || c == '\r' })
	for _, val := range strArr {
		reqArr := strings.Fields(val)
		
		Id, _ := strconv.Atoi(reqArr[0])
		UserRole, _ := strconv.Atoi(reqArr[4])
		
		bk := UserFromDB{}
		bk = UserFromDB{Id: Id, UserName: reqArr[1], Login: reqArr[2], Pswd: reqArr[3], UserRole: UserRole}
		//bk = UserFromDB{Id: 0, UserName: "Admin", Login: "Admin", UserRole: 5}
		fmt.Println(Id, UserRole, bk)
		bks.Users = append(bks.Users, bk)
		
	}
	
	fmt.Println(bks.Users)

	// Response template
	tmpl, _ := template.ParseFiles("www/shell.html")
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, bks)
}


