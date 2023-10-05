package main

import (
	//"fmt"
	"net/http"
	"strings"
	"strconv"
	//"io/ioutil"
	"github.com/gorilla/sessions"
	//"bufio"
	//"regexp"
)

// For Login
var cookieStore = sessions.NewCookieStore([]byte("secret"))

const cookieName = "MyCookie"

type sesKey int

type UserFromDB struct {
    Id int
    UserName string
    Login string
    Pswd string
    UserRole int
}

const (
	sesKeyLogin sesKey = iota

)

// Users handler
func login(w http.ResponseWriter, r *http.Request) {
	// parameters from POST
	backLink := r.FormValue("backLink")

	log := r.FormValue("login")
	pass := r.FormValue("password")

	//fmt.Println(log)
	//fmt.Println(pass)

	// open database
	//db := openDB()
	//defer db.Close() // close database

	// select query
	//rows, err := db.Query(`SELECT * FROM users`)
	//CheckError(err)
	//defer rows.Close()

	//for rows.Next() {
		openDB("http://192.168.63.60", "8086") // "http://192.168.1.136"  "http://192.168.63.60"
		req := db_request("select Id, UserName, Login, Pswd, UserRole from users;") 

		strArr := strings.FieldsFunc(req, func(c rune) bool { return c == '\n' || c == '\r' })
		for _, val := range strArr {
			reqArr := strings.Fields(val)
		
		Id, _ := strconv.Atoi(reqArr[0])
		UserRole, _ := strconv.Atoi(reqArr[4])
		bk := UserFromDB{Id: Id, UserName: reqArr[1], Login: reqArr[2], Pswd: reqArr[3], UserRole: UserRole}

		if string(bk.Login) == string(log) && bk.Pswd == pass {
			//fmt.Println("OK")
			ses, err := cookieStore.Get(r, cookieName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			ses.Values[sesKeyLogin] = log
			ses.Values["sesUserRole"] = bk.UserRole
			err = cookieStore.Save(r, w, ses)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			break
		}
	}
	//}
	//CheckError(err)
	http.Redirect(w, r, backLink, http.StatusSeeOther)

}

// LogOut handler
func logout(w http.ResponseWriter, r *http.Request) {
	backLink := r.FormValue("backLink")

	for _, v := range r.Cookies() {
		c := http.Cookie{
			Name:   v.Name,
			MaxAge: -1}
		http.SetCookie(w, &c)
	}
	http.Redirect(w, r, backLink, http.StatusSeeOther)
}

// Check Cookies
func check_cookies(w http.ResponseWriter, r *http.Request) (string, int){
	var userRole int = 0
	ses, err := cookieStore.Get(r, cookieName)
	CheckErrorHttp(err, w)

	log, ok := ses.Values[sesKeyLogin].(string)
	if !ok {
		log = ""
	} else {
		userRole, _ = ses.Values["sesUserRole"].(int)
	}
	
	return log, userRole
}




