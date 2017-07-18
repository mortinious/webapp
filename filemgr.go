package main

import (
	//"fmt"
	"bufio"
	"os"
	"strings"
)

type Obj struct{
	Vars map[string]string
}

func loadListFile(fn string){
	dat, err := os.Open("/data/"+fn)
	if(err != nil){
		panic(err)
	}
	defer dat.Close()
	
	scan := bufio.NewScanner(dat)
	for scan.Scan(){
		
	}
}

func readObject(s string) Obj{
	obj := Obj{Vars: map[string]string{}}
	list1 := strings.Split(s, ";")
	for i := range list1 {
		s2 := strings.Split(list1[i], ":")
		obj.Vars[s2[0]] = s2[1]
	}
	return obj
}

