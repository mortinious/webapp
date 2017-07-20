/*
This file manages http requests and html templates
*/

package main

import(
	"net/http"
	"html/template"
	"strconv"
	"fmt"
)

//Handles requests from root ("/")
func handleRoot(w http.ResponseWriter, r *http.Request) {
	//Gets session id from cookie or creates a new if no exists (from sessionmgr.go)
	sid := initSession(w,r)
	
	//Init messages to be displayed on page
	messages := make(map[string]string)
	
	//Checks clients cookie is old session exists
	isNew := newSession(strconv.Itoa(sid))
	
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
		username, success := getUsernameBySessionId(strconv.Itoa(sid))
		if(success){ //If session was found
			if(username != ""){ //If user exists
			
				//Get userinfo from userinfo table. Load into new User struct
				user, success := getUserByUsername(username)
				if(success){ //If user was found in userinfo
				
					//Sets username and loads mypage.html into index
					messages["username"] = user.name
					index, _ = template.ParseFiles("html/mypage.html")
				}
			}else{ //If no user was found from sessionid
				
				index, _ = template.ParseFiles("html/login.html")
				messages["no_user_found"] = ""
			}
		}
	}
	
	//Sends modified html to client
	head.Execute(w, nil)
	index.Execute(w, messages) //Inserts messages into index template
	foot.Execute(w, nil)
	
	//Sends a string2 of text for debugging
	fmt.Fprintf(w , "Server running on: "+dockerid+"<br>Session ID: "+strconv.Itoa(sid))
}

//Handles login requests from /login
func handleLogin(w http.ResponseWriter, r *http.Request) {
	//Continue only if request method is POST
	if(r.Method != "POST"){return}
	
	//Get session id from cookie
	sid := strconv.Itoa(initSession(w,r))
	
	
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
	}
	
	//Redirect back to root
	http.Redirect(w, r, "/", 302)
}

//Handles register requests
func handleRegister(w http.ResponseWriter, r *http.Request) {
	//Continue only if request method is POST
	if(r.Method != "POST"){return}
	
	//Get session id from cookie
	sid := strconv.Itoa(initSession(w,r))
	
	//Load variables from html form from login.html and pass it to vars
	r.ParseForm()
	user := r.Form.Get("user")
	pass := r.Form.Get("pass")
	
	//Write user credentials to sql
	addUser(user,pass)
	
	//Update session db with user
	setUserForSession(sid,user)
	
	//Redirect back to root
	http.Redirect(w, r, "/", 302)
}

//Handles Logout requests
func handleLogout(w http.ResponseWriter, r *http.Request) {
	//Continue only if request method is POST
	if(r.Method != "POST"){return}
	
	//Get session id from cookie
	sid := strconv.Itoa(initSession(w,r))
	
	//Remove user from sessionid in sessions-table
	setUserForSession(sid,"")
	
	//Redirect back to root
	http.Redirect(w, r, "/", 302)
}

//Handles requests for favicon.ico to take these requests away from root
func handleIcon(w http.ResponseWriter, r *http.Request) {
	//Code for favicon goes here!
}

