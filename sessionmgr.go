/*
This file managed cookies with session data
*/

package main

import(
	"fmt"
	"net/http"
	"github.com/gorilla/sessions"
	"math/rand"
)

//Creates new cookie store
var store = sessions.NewCookieStore([]byte("verysecretpasscode"))

//Reads sessionid from cookie or creates new on if no is found
func initSession(w http.ResponseWriter, r *http.Request) int{
	//Updates seed with random strgin from last seed
	rand.Seed(rand.Int63())
	
	//Get session, or create one if no exists from name
	session, err := store.Get(r, "session")
	
	//if error message is returned, display it
	if(err != nil){
		fmt.Println("Error: no session found"+err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return -1
	}
	
	//Inits sid var
	sid := -1
	
	//Reads or sets sessionid
	if(session.Values["ID"] != nil){
		sid = session.Values["ID"].(int)
	}else{
		sid = rand.Int()
	}
	
	//Updates cookie to keep it updated
	session.Values["ID"] = sid
	session.Save(r, w)
	
	fmt.Println("New session created with ID:",sid)
	
	//Returns session id
	return sid
}