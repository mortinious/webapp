/*
This file manages http requests and html templates
*/

package main

import(
	"github.com/gorilla/mux"
	"net/http"
	"html/template"
	"fmt"
)

func routeTrafic(){
	r := mux.NewRouter()
	r.HandleFunc("/", handleRoot)
	r.HandleFunc("/login", handleLogin).Methods("POST")
	r.HandleFunc("/logout", handleLogout).Methods("POST")
	r.HandleFunc("/reg", handleRegister).Methods("POST")
	http.Handle("/", r)
}

//Handles requests from root ("/")
func handleRoot(w http.ResponseWriter, r *http.Request) {
	//Gets session id from cookie or creates a new if no exists (from sessionmgr.go)
	sid := initSession(w,r)
	
	//Init messages to be displayed on page
	messages := make(map[string]string)
	
	//TESTING GROUNDS
	fmt.Println("Header redirection message: "+w.Header().Get("message"))
	
	//Checks clients cookie is old session exists
	isNew := newSession(sid)
	
	//Reads html-templates into vars
	head, _ := template.ParseFiles("html/header.html")
	index, _ := template.ParseFiles("html/index.html")
	foot, _ := template.ParseFiles("html/footer.html")
	
	//Autologin
	if(isNew){//No cookie found and a new session was created
		
		fmt.Println("New session created")
		
		//Sets index to login page as new sessions don't have a user appended
		index, _ = template.ParseFiles("html/login.html")
	}else{//Cookie found and session id loaded
		
		//gets username from databse based on sessionid
		username, success := getUsernameBySessionId(sid)
		if(success){ //If session was found
			if(username != ""){ //If user exists
				//Sets username and loads mypage.html into index
				messages["username"] = username
				index, _ = template.ParseFiles("html/mypage.html")
			}else{ //If no user was found from sessionid
				
				index, _ = template.ParseFiles("html/login.html")
				messages["no_user_found"] = getSessionData(sid, "login_status")
			}
		}
	}
	
	//Sends modified html to client
	head.Execute(w, nil)
	index.Execute(w, messages) //Inserts messages into index template
	foot.Execute(w, nil)
	
	//Sends a string2 of text for debugging
	fmt.Fprintf(w , "Server running on: "+dockerid+"<br>Session ID: "+sid)
	
	messages = make(map[string]string)
}

//Handles login requests from /login
func handleLogin(w http.ResponseWriter, r *http.Request) {
	//Get session id from cookie
	sid := initSession(w,r)
	
	
	r.ParseForm()
	user := r.Form.Get("user")
	pass := r.Form.Get("pass")
	
	//Check with sql is credentials is correct
	success := checkUserCredentials(user,pass)
	if(success){ //If correct
		fmt.Println("Login success")
		
		//Update session-table with new user for corresponding sessionid
		setUserForSession(sid,user)
		
		
	}else{
		//TODO: Send failed login message
		setSessionDataWithDataLife(sid, "login_status", "Incorrect password or username", 1)
	}
	
	//Redirect back to root
	http.Redirect(w, r, "/", 302)
}

//Handles register requests
func handleRegister(w http.ResponseWriter, r *http.Request) {
	//Continue only if request method is POST
	if(r.Method != "POST"){return}
	
	//Get session id from cookie
	sid := initSession(w,r)
	
	//Load variables from html form from login.html and pass it to vars
	r.ParseForm()
	user := r.Form.Get("user")
	pass := r.Form.Get("pass")
	
	//Check if user don't exists
	if(!checkUsername(user)){
	
		//Write user credentials to sql
		addUser(user,pass)
		
		setSessionDataWithDataLife(sid, "login_status", "New user account created.", 1)

	}else{
		setSessionDataWithDataLife(sid, "login_status", "User already exists. Try another name.", 1)
	}
	//Redirect back to root
	http.Redirect(w, r, "/", 302)
}

//Handles Logout requests
func handleLogout(w http.ResponseWriter, r *http.Request) {
	//Continue only if request method is POST
	if(r.Method != "POST"){return}
	
	//Get session id from cookie
	sid := initSession(w,r)
	
	//Remove user from sessionid in sessions-table
	setUserForSession(sid,"")
	
	//Redirect back to root
	http.Redirect(w, r, "/", 302)
}

//Handles requests for favicon.ico to take these requests away from root
func handleIcon(w http.ResponseWriter, r *http.Request) {
	//Code for favicon goes here!
}

