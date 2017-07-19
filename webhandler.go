package main

import(
	"net/http"
	"io/ioutil"
	"strconv"
	"fmt"
)

type Page struct{
	Title string
	Body []byte
}

func loadPage (fn string) *Page{
	path := "html/index.html"
	if(fn == ""){
		path = "html/index.html"
	}else{
		path = "html/"+fn
	}
	body, _ := ioutil.ReadFile(path)
	return &Page{"Test", body}
}

//Handles requests from root ("/")
func handleRoot(w http.ResponseWriter, r *http.Request) {
	sid := initSession(w,r)
	username := "Not logged in"
	isNew := newSession(strconv.Itoa(sid))
	if(isNew){
		//New session was created
		fmt.Println("New session created")
		
	}else{
		user, success := getUserBySessionId(strconv.Itoa(sid))
		if(success){
			username = user.name
		}
	}
	num++
	fmt.Fprintf(w , "%sThis is a Go webserver! <br>This webpage has been loaded "+strconv.Itoa(num)+" times.<br>Server running on: "+dockerid+"<br>Session ID: "+strconv.Itoa(sid)+"<br>You are logged in as: "+username, loadPage("").Body)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login request")
	sid := strconv.Itoa(initSession(w,r))
	if(r.Method != "POST"){return}
	r.ParseForm()
	user := r.Form.Get("user")
	pass := r.Form.Get("pass")
	addUser(user,pass)
	setUserForSession(sid,user)
	
	http.Redirect(w, r, "/", 302)
}

//Handles requests for favicon.ico to take these requests away from root
func handleIcon(w http.ResponseWriter, r *http.Request) {
	//Code for favicon goes here!
}

