/*
This is the main file
*/

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



//Main function
func main(){
	//Sets port number and runs init()
	fmt.Println("Starting server...")
	port := 7001
	
	//Init global vars
	initVars()

	//Creates new seed based on prevoius seed
	rand.Seed(rand.Int63())
	
	//Registers handlers for GET and POST requests
	http.HandleFunc("/favicon.ico", handleIcon)
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	
	//Prints start message
	fmt.Println("Server Started on port:"+strconv.Itoa(port))
	
	//Initializes listening thread
	http.ListenAndServe(":"+strconv.Itoa(port), context.ClearHandler(http.DefaultServeMux))
}

