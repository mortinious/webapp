package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"bufio"
)

var mysqlConnectArgs string

type User struct{
	name string
	pass string
}

func mysqlConnection() string{
	file, err := os.Open("conf/mysql.conf")
	if(err != nil){
		fmt.Println("MySQL file not found")
		return "webapp:password@/webapp?charset=utf8"
	}
	s := bufio.NewScanner(file)
	
	s.Scan()
	user := s.Text()
	s.Scan()
	pass := s.Text()
	s.Scan()
	addr := s.Text()
	s.Scan()
	tabl := s.Text()
	s.Scan()
	args := s.Text()
	
	

	return user+":"+pass+"@"+addr+"/"+tabl+"?"+args
}

func check(e error){
	if(e != nil){
		fmt.Println(e.Error())
	}
}

func addUser(name string, pass string) bool{
	if(name == "" || pass == ""){return false}
	
	fmt.Println("Attempting to add user with name: "+name)
	
	db, err := sql.Open("mysql", mysqlConnection())
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
	
	db, err := sql.Open("mysql", mysqlConnection())
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
	fmt.Println("Attempting to set User of session id: "+session+" to: "+username)
	
	
	
	db, err := sql.Open("mysql", mysqlConnection())
	check(err)
	
	defer db.Close()
	
	_, err = db.Query("UPDATE sessions SET user='"+username+"' WHERE session='"+session+"'")
	check(err)
}

func checkUserCredentials(username string, password string) bool{
	fmt.Println("Attempting to check credentials for: "+username)
	if(username == "" || password == ""){return false}
	
	var user, pass string
	
	db, err := sql.Open("mysql", mysqlConnection())
	check(err)
	
	defer db.Close()
	
	rows, err := db.Query("SELECT * FROM userinfo WHERE (username='"+username+"' AND password='"+password+"')")
	check(err)
	
	if(rows != nil){
		for rows.Next(){
			err = rows.Scan(&user, &pass)
			if(user == username && pass == password){
				return true
			}
		}
	}
	return false
}


func getUsernameBySessionId(session string) (string, bool){
	if(session == ""){return "", false}
	fmt.Println("Attempting to get User by session id: "+session)
	
	var username string
	
	db, err := sql.Open("mysql", mysqlConnection())
	check(err)
	
	defer db.Close()
	
	rows, err := db.Query("SELECT user FROM sessions WHERE session='"+session+"'")
	check(err)
	
	if(rows != nil){
		rows.Next()
		err = rows.Scan(&username)
		check(err)
		fmt.Println("-- Got: "+username)
		rows.Close()
		return username, true
	}
	return "", false
}

func getUserByUsername(username string) (User, bool){
	fmt.Println("Attempting to get User object from username: "+username)
	if(username == ""){return User{name: "", pass: ""}, false}
	
	var user, pass string
	
	db, err := sql.Open("mysql", mysqlConnection())
	check(err)
	
	defer db.Close()
	
	rows, err := db.Query("SELECT * FROM userinfo WHERE username='"+username+"'")
	check(err)
	
	if(rows != nil){
		for rows.Next(){
			err = rows.Scan(&user, &pass)
			if(user == username){
				return User{name: user, pass: pass}, true
			}
		}
	}
	return User{name: "", pass: ""}, false
}
