package main

import(
	"net/http"
	"html/template"
	"strconv"
	"fmt"
)

//Handles requests from root ("/")
func handleRoot(w http.ResponseWriter, r *http.Request) {
	sid := initSession(w,r)
	messages := make(map[string]string)
	
	isNew := newSession(strconv.Itoa(sid))
	
	head, _ := template.ParseFiles("html/header.html")
	index, _ := template.ParseFiles("html/index.html")
	foot, _ := template.ParseFiles("html/footer.html")
	
	//Check if sessionid is exists
	if(isNew){
		//New session was created
		fmt.Println("New session created")
		index, _ = template.ParseFiles("html/login.html")
	}else{
		username, success := getUsernameBySessionId(strconv.Itoa(sid))
		if(success){
			if(username != ""){
				user, success := getUserByUsername(username)
				if(success){
					messages["username"] = user.name
					index, _ = template.ParseFiles("html/mypage.html")
				}
			}else{
				
				index, _ = template.ParseFiles("html/login.html")
				messages["no_user_found"] = ""
			}
		}
	}
	num++
	
	head.Execute(w, nil)
	index.Execute(w, messages)
	foot.Execute(w, nil)
	fmt.Fprintf(w , "This is a Go webserver! <br>This webpage has been loaded "+strconv.Itoa(num)+" times.<br>Server running on: "+dockerid+"<br>Session ID: "+strconv.Itoa(sid))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login Request")
	sid := strconv.Itoa(initSession(w,r))
	if(r.Method != "POST"){return}
	r.ParseForm()
	user := r.Form.Get("user")
	pass := r.Form.Get("pass")
	success := checkUserCredentials(user,pass)
	if(success){
		fmt.Println("Login success")
		setUserForSession(sid,user)
	}else{
		//Send failed login message
	}
	
	http.Redirect(w, r, "/", 302)
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	sid := strconv.Itoa(initSession(w,r))
	if(r.Method != "POST"){return}
	r.ParseForm()
	user := r.Form.Get("user")
	pass := r.Form.Get("pass")
	addUser(user,pass)
	setUserForSession(sid,user)
	
	http.Redirect(w, r, "/", 302)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	sid := strconv.Itoa(initSession(w,r))
	if(r.Method != "POST"){return}
	r.ParseForm()
	setUserForSession(sid,"")
	
	http.Redirect(w, r, "/", 302)
}

//Handles requests for favicon.ico to take these requests away from root
func handleIcon(w http.ResponseWriter, r *http.Request) {
	//Code for favicon goes here!
}

