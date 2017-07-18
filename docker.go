package main

import(
	"os"
)

//Reads env.variable "DOCKER_ID" for the dockerID and returns the value, if the variable do not exixts it returns false, otherwise true
func getDockerID() (string){
	id, _ := os.Hostname()
	return id
}

