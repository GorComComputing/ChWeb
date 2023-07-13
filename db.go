package main

import (
	//"fmt"
	"net/http"
	"io/ioutil"
	"strings"
)


var db_host string 
var db_port string 


// open database
func openDB(host string, port string) {
	db_host = host
	db_port = port
}


// SQL-запрос к базе данных 
func db_request(db_Query string) string {
	var output string
	
	db_Query = strings.Replace(db_Query, " ", "%20", -1)
	db_Query = strings.Replace(db_Query, ";", "%3b", -1)
	resp, err := http.Get(db_host + ":" + db_port + "/api?cmd=" + db_Query)
	
	if err != nil {
		output = "Request FAIL: " + err.Error() + "\n"
		return output
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		output = "Request FAIL: " + err.Error() + "\n"
		return output
	}

	return string(body) 
}
