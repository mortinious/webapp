//This is dev only
package main

import(
	"net/http"
	"fmt"
	"strconv"
	"math/rand"
	
	"github.com/gorilla/context"
)

//Init global variables
var num = 0
var dockerid = "NO_HOSTNAME_FOUND"

func initVars(){
	//Gets docker ID from local machine
	dockerid = getDockerID()
}

//Handles requests from root ("/")
func handleRoot(w http.ResponseWriter, r *http.Request) {
	sid := initSession(w,r)
	num++
	fmt.Fprintf(w , "This is a Go webserver! \nThis webpage has been loaded "+strconv.Itoa(num)+" times.\nServer running on: "+dockerid+"\nSession ID: "+strconv.Itoa(sid))
}

//Handles requests for favicon.ico to take these requests away from root
func handleIcon(w http.ResponseWriter, r *http.Request) {
	//Code for favicon goes here!
}

//Main function
func main(){
	//Sets port number and runs init()
	port := 7001
	initVars()
	initRadix()

	//Creates new seed based on prevoius seed
	rand.Seed(rand.Int63())
	
	//Test redis connection
	connect("192.168.42.50",6379)
	
	
	//Registers handlers for GET and POST requests
	http.HandleFunc("/favicon.ico", handleIcon)
	http.HandleFunc("/", handleRoot)
	
	//Prints start message
	fmt.Println("Server Started on port:"+strconv.Itoa(port))
	
	//Initializes listening thread
	http.ListenAndServe(":"+strconv.Itoa(port), context.ClearHandler(http.DefaultServeMux))
}

