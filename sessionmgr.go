/*8
This file managed cookies with session data
*/

package main

import(
	"fmt"
	"net/http"
	"github.com/gorilla/sessions"
	"math/rand"
	"strconv"
)

//Define sessiondata struct
type Data struct{
	data string
	life int
}

//Defines sessiondata var, USAGE: sessiondata["sessionid"]["key"]
var sessiondata = make(map[string]map[string]Data)

//Creates new cookie store
var store = sessions.NewCookieStore([]byte("verysecretpasscode"))

func setSessionData(sid string, key string, val string){
	if(sessiondata[sid] == nil){
		sessiondata[sid] = make(map[string]Data)
	}
	sessiondata[sid][key] = Data{data: val, life: 0}
}

func setSessionDataWithDataLife(sid string, key string, val string, life int){
	if(sessiondata[sid] == nil){
		sessiondata[sid] = make(map[string]Data)
	}
	sessiondata[sid][key] = Data{data: val, life: life}
}

func getSessionData(sid string, key string) string{
	data := sessiondata[sid][key]
	if(data.data != ""){
		if(data.life > 0){
			data.life -= 1
			if(data.life == 0){
				delete(sessiondata[sid], key)
				return data.data
			}
			setSessionDataWithDataLife(sid, key, data.data, data.life)
			return data.data
		}
		return data.data
	}
	return ""
}

//Reads sessionid from cookie or creates new on if no is found
func initSession(w http.ResponseWriter, r *http.Request) string{
	//Updates seed with random strgin from last seed
	rand.Seed(rand.Int63())
	
	//Get session, or create one if no exists from name
	session, err := store.Get(r, "session")
	
	//if error message is returned, display it
	if(err != nil){
		fmt.Println("Error: no session found"+err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "-1"
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
	return strconv.Itoa(sid)
}