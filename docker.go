package main

import(
	"os"
)

func getDockerID() (string){
	id, _ := os.Hostname()
	return id
}

