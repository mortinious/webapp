package main

import (
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	"strconv"
	"os"
	"bufio"
)

var rediskey = ""
var hostname = ""

func initRadix(){
	rediskey, hostname = getRedisConf()
}

func getRedisConf() (string,string){
	fn , err := os.Open("redis.conf")
	if(err != nil){return ""}
	
	scan := bufio.NewScanner(fn)
	scan.Scan()
	key := scan.Text()
	scan.Scan()
	return scan.Text(), key 
}

func connect(){
	p, err := pool.New("tcp", hostname, 10)
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
