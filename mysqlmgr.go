package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	//"strings"
)

type User struct{
	name string
	pass string
}

func check(e error){
	if(e != nil){
		fmt.Println(e.Error())
	}
}

func addUser(name string, pass string) bool{
	if(name == "" || pass == ""){return false}
	
	fmt.Println("Attempting to add user with name: "+name)
	
	db, err := sql.Open("mysql", "webapp:password@/webapp?charset=utf8")
	check(err)
	
	defer db.Close()
	
	_, err = db.Query("INSERT INTO userinfo (username, password) VALUES ('"+name+"', '"+pass+"')")
	if(err != nil){
		fmt.Println(err.Error())
		return false
	}else{
		fmt.Println("User added "+name)
		return true
	}
}

func newSession(session string) bool{
	fmt.Println("Initiating new session")
	if(session == ""){return false}
	
	db, err := sql.Open("mysql", "webapp:password@/webapp?charset=utf8")
	check(err)
	
	defer db.Close()
	
	_, err = db.Query("INSERT INTO sessions (session) VALUES ('"+session+"')")
	check(err)
	if(err == nil){
		fmt.Println("New session established, id: "+session)
		return true
	}
	fmt.Println("Existing session established, id: "+session)
	return false
}

func setUserForSession(session string, username string){
	if(session == ""){return}
	fmt.Println("Attempting to set User by session id: "+username)
	
	
	
	db, err := sql.Open("mysql", "webapp:password@/webapp?charset=utf8")
	check(err)
	
	defer db.Close()
	
	_, err = db.Query("UPDATE sessions SET user='"+username+"' WHERE session='"+session+"'")
	check(err)
}


func getUserBySessionId(session string) (User, bool){
	if(session == ""){return User{"",""}, false}
	fmt.Println("Attempting to get User by session id: "+session)
	
	var username, pass string
	
	db, err := sql.Open("mysql", "webapp:password@/webapp?charset=utf8")
	check(err)
	
	defer db.Close()
	
	rows, err := db.Query("SELECT user FROM sessions WHERE session='"+session+"'")
	check(err)
	
	if(rows != nil){
		rows.Next()
		err = rows.Scan(&username)
		check(err)
		
		rows, err := db.Query("SELECT * FROM userinfo WHERE username='"+username+"'")
		check(err)
		if(rows != nil){
			rows.Next()
			err = rows.Scan(&username, &pass)
			check(err)
			return User{username, pass}, true
		}
		fmt.Println(rows.Err())
	}
	return User{"",""}, false
}
