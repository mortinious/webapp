/*
This file managed functions related to docker
*/

package main

import(
	"os"
)

//Reads hostname from OS
func getDockerID() (string){
	id, _ := os.Hostname()
	return id
}

