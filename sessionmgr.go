package main

import(
	"fmt"
	"net/http"
	"github.com/gorilla/sessions"
	"math/rand"
)

var store = sessions.NewCookieStore([]byte("verysecretpasscode"))

func initSession(w http.ResponseWriter, r *http.Request) int{
	rand.Seed(rand.Int63())
	//Get session, or create one if no exists from name
	session, err := store.Get(r, "session")
	
	//if error message is returned, display it
	if(err != nil){
		fmt.Println("Error: no session found"+err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return -1
	}
	
	sid := -1
	
	if(session.Values["ID"] != nil){
		sid = session.Values["ID"].(int)
	}else{
		sid = rand.Int()
	}
	
	session.Values["ID"] = sid
	session.Save(r, w)
	fmt.Println("New session created with ID:",sid)
	return sid
}