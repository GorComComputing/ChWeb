package main

import (
	"fmt"
	"html/template"
	"net/http"
	
    	"os/exec"
    	"os"
 	"io"
    	"bufio"
    	"log"
    	"strings"
)

// Структура config-файла
type Config struct {
	Leapsectz	string
	Driftfile    	string
	Makestep       	string
	Makestep2      	string
	Rtcsync 	bool
	Logdir       	string
	Local    	string
	Server		string
	Allow 		string
}

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


// Shell handler
func shell(w http.ResponseWriter, r *http.Request) {
	
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


// Запуск Chrony
func start(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start")

	cmd := exec.Command("/etc/init.d/chrony", "start")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}
	fmt.Println(string(out))

	//fmt.Fprintf(w, "<p class='bg-success table-page'>Chrony запущен</p>")
	messages <- string("Chrony запущен")
	isStartChrony()
}


// Остановка Chrony
func stop(w http.ResponseWriter, r *http.Request) {
	fmt.Println("stop")
	
	cmd := exec.Command("/etc/init.d/chrony", "stop")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}
	fmt.Println(string(out))

	messages <- string("Chrony остановлен")
	isStartChrony()
}


// Перезапуск Chrony
func restart(w http.ResponseWriter, r *http.Request) {

	cmd := exec.Command("/etc/init.d/chrony", "restart")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}
	fmt.Println(string(out))

	messages <- string("Chrony перезапущен")
	isStartChrony()
}


// Проверяет запущен ли Chrony
func isStartChrony() bool {
	file, err := os.Open("/var/run/chrony/chronyd.pid")
	defer file.Close()

	//handle errors while opening
	if err != nil {
		messages2 <- string("<svg xmlns='http://www.w3.org/2000/svg' width='16' height='16' fill='#dc3545' class='bi bi-circle-fill' viewBox='0 0 16 16'><circle cx='8' cy='8' r='8'/></svg>")
		return false
	} else {
		messages2 <- string("<svg xmlns='http://www.w3.org/2000/svg' width='16' height='16' fill='#198754' class='bi bi-circle-fill' viewBox='0 0 16 16'><circle cx='8' cy='8' r='8'/></svg>")
		return true
	}
}


// Читает config-файл
func scan() (Config, string) {

	var config Config

	// open the file
	file, err := os.Open("/etc/pzg-chrony.conf")
	
	//handle errors while opening
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	fileScanner := bufio.NewScanner(file)
	
	// read line by line
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if len(line) > 0 && string(line[0]) != "#" {
			words := strings.Fields(fileScanner.Text())
		
			switch words[0] {
   			case "leapsectz": config.Leapsectz = words[1]
			case "driftfile": config.Driftfile = words[1]
			case "makestep": {config.Makestep = words[1]; config.Makestep2 = words[2]}
			case "rtcsync": config.Rtcsync = true
			case "logdir": config.Logdir = words[1]
			case "local": config.Local = words[2]
			case "server": config.Server = words[1]
			case "allow": config.Allow = words[1]
		
			default: fmt.Println("Unknown directive")
			}
		}
	}
	
	fmt.Println(config)
	
	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	defer file.Close()
	
	dat, err := os.ReadFile("/etc/pzg-chrony.conf")
    	fmt.Print(string(dat))
	
	return config, string(dat)
}


// Сохраняет config-файл и перезапускает Chrony
func save(w http.ResponseWriter, r *http.Request) {
	// parameters from POST
	Leapsectz := r.FormValue("Leapsectz")
	Driftfile := r.FormValue("Driftfile")
	Makestep := r.FormValue("Makestep")
	Makestep2 := r.FormValue("Makestep2")
	Rtcsync := r.FormValue("Rtcsync")
	fmt.Println(Rtcsync)
	Logdir := r.FormValue("Logdir")
	Local := r.FormValue("LocalStratum")
	Server := r.FormValue("Server")
	Allow := r.FormValue("Allow")
	
    	  	
    	file, err := os.OpenFile("./files/tmp.conf", os.O_TRUNC | os.O_WRONLY, 0600)
    	if err != nil {
        	fmt.Println("Unable to open file:", err) 
        	os.Exit(1) 
    	}
    	defer file.Close()

	
	if Leapsectz != "" { 
    	if _, err = file.WriteString("leapsectz " + Leapsectz + "\n"); err != nil {
    		fmt.Println("Unable to write string:", err) 
        	os.Exit(1) 
    	}
    	}
    	if Driftfile != "" { 
    	if _, err = file.WriteString("driftfile " + Driftfile + "\n"); err != nil {
    		fmt.Println("Unable to write string:", err) 
        	os.Exit(1) 
    	}
    	}
    	if Makestep != "" && Makestep2 != "" { 
    	if _, err = file.WriteString("makestep " + Makestep + " " + Makestep2 + "\n"); err != nil {
    		fmt.Println("Unable to write string:", err) 
        	os.Exit(1) 
    	}
    	}
    	if Rtcsync != "false" { 
    	if _, err = file.WriteString("rtcsync\n"); err != nil {
    		fmt.Println("Unable to write string:", err) 
        	os.Exit(1) 
    	}
    	}
    	if Logdir != "" { 
    	if _, err = file.WriteString("logdir " + Logdir + "\n"); err != nil {
    		fmt.Println("Unable to write string:", err) 
        	os.Exit(1) 
    	}
    	}
    	if Local != "" { 
    	if _, err = file.WriteString("local stratum " + Local + "\n"); err != nil {
    		fmt.Println("Unable to write string:", err) 
        	os.Exit(1) 
    	}
    	}
    	if Server != "" { 
    	if _, err = file.WriteString("server " + Server + "\n"); err != nil {
    		fmt.Println("Unable to write string:", err) 
        	os.Exit(1) 
    	}
    	}
    	if Allow != "" { 
    	if _, err = file.WriteString("allow " + Allow + "\n"); err != nil {
    		fmt.Println("Unable to write string:", err) 
        	os.Exit(1) 
    	}
    	}
    	
    	// перенести из tmp в основной файл
    	cmd := exec.Command("cp", "./files/tmp.conf", "/etc/pzg-chrony.conf")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not back copy: ", err)
	}
	fmt.Println(string(out))
    	fmt.Println("Saved OK")
    	restart(w, r)
    	
    	
    	_ , File := scan()
	//File = strings.Replace(File, "\n", "<br/>", -1)
	//bks.File = template.HTML(File)
    	
    	fmt.Fprintf(w, File)
    	messages <- string("Config-файл Chrony сохранен<br/>Chrony запущен")
}


// Проверяет активность Chrony
func activity(w http.ResponseWriter, r *http.Request) {
	
	cmd := exec.Command("chronyc", "activity")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}
	fmt.Println(string(out))

	fmt.Fprintf(w, string(out))
}


// chronyc tracking - информация о сервере времени
func tracking(w http.ResponseWriter, r *http.Request) {
	
	cmd := exec.Command("chronyc", "tracking")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}
	fmt.Println(string(out))

	fmt.Fprintf(w, string(out))
}


// chronyc sources -v   список источников времени
func sources(w http.ResponseWriter, r *http.Request) {
	
	cmd := exec.Command("chronyc", "sources", "-v")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}
	fmt.Println(string(out))

	fmt.Fprintf(w, string(out))
}


// chronyc sourcestats -v список источников времени
func sourcestats(w http.ResponseWriter, r *http.Request) {
	
	cmd := exec.Command("chronyc", "sourcestats", "-v")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}
	fmt.Println(string(out))

	fmt.Fprintf(w, string(out))
}


// chronyc clients - список подключенных клиентов
func clients(w http.ResponseWriter, r *http.Request) {
	
	cmd := exec.Command("chronyc", "clients")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}
	fmt.Println(string(out))

	fmt.Fprintf(w, string(out))
}


// Читает Config-файл
func config(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Config")
    	_ , File := scan()
    	fmt.Fprintf(w, File)
    	messages <- string("Config-файл прочитан")
}


// Сохраняет Config-файл
func saveconfig(w http.ResponseWriter, r *http.Request) {
	// parameters from POST
	textarea := r.FormValue("file")
	
	fmt.Println("save config")
	fmt.Println(textarea)
		
    	  	
    	file, err := os.OpenFile("./files/tmp.conf", os.O_TRUNC | os.O_WRONLY, 0600)
    	if err != nil {
        	fmt.Println("Unable to open file:", err) 
        	os.Exit(1) 
    	}
    	defer file.Close()

	
	if textarea != "" { 
    	if _, err = file.WriteString(textarea); err != nil {
    		fmt.Println("Unable to write string:", err) 
        	os.Exit(1) 
    	}
    	}
    	
    	
    	// перенести из tmp в основной файл
    	cmd := exec.Command("cp", "./files/tmp.conf", "/etc/pzg-chrony.conf")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not back copy: ", err)
	}
	fmt.Println(string(out))
    	fmt.Println("Saved OK")
    	restart(w, r)
    	
    	_ , File := scan()
    	
    	fmt.Fprintf(w, File)
    	messages <- string("Config-файл Chrony сохранен<br/>Chrony запущен")
}


// Восстановление config-файла
func restore(w http.ResponseWriter, r *http.Request) {
	
	fmt.Println("restore config") 	
    	
    	// перенести из files в основной файл
    	cmd := exec.Command("cp", "./files/pzg-chrony.conf", "/etc/pzg-chrony.conf")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not back copy: ", err)
	}
	fmt.Println(string(out))
    	fmt.Println("Restored OK")
    	restart(w, r)
    	
    	_ , File := scan()
    	
    	fmt.Fprintf(w, File)
    	messages <- string("Config-файл Chrony восстановлен<br/>Chrony запущен")
}



