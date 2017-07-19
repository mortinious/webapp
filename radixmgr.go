package main

import (
	"fmt"
//	"github.com/mediocregopher/radix.v2/pool" //Removed Radix for testing purpose
	"strconv"
	"os"
	"bufio"
)

var rediskey = ""

func initRadix(){
	rediskey = getRedisKey()
}

func getRedisKey() string{
	fn , err := os.Open("rediskey")
	if(err != nil){return ""}
	
	scan := bufio.NewScanner(fn)
	scan.Scan()
	return scan.Text()
}

func connect(url string, port int){
	p, err := pool.New("tcp", url+":"+strconv.Itoa(port), 10)
	if(err != nil){
		fmt.Println(err.Error())
		return
	}
	
	con, err := p.Get()
	if(err != nil){
		fmt.Println(err.Error())
		return
	}
	
	defer p.Put(con)
	
	err = con.Cmd("AUTH", rediskey).Err
	if(err != nil){
		fmt.Println(err.Error())
		return
	}
	
	err = con.Cmd("SET", "Entry", "This is something").Err
	if(err != nil){
		fmt.Println(err.Error())
		return
	}
	
	s, err := con.Cmd("GET", "Entry").Str()
	if(err != nil){
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Redis get query: "+s)
	

}

func set(key string, val string){
	
	
}
