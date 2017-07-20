/*
This file managed mysql connections and queries

This part needs a mysql server with the collowing config:

- Databse with 2 tables:
--- userinfo with 2 columns
----- username (PRIMARY)
----- password (non-null)
--- sessions with 2 columns
----- session (PRIMARY)
----- user (default NULL)

- User created with SELECT, INSERT and UPDATE permissions for the table

*/

package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"bufio"
)

//Declares User struct
type User struct{
	name string
	pass string
}

//Get mysql connection args from mysql.conf
//TODO: Update config file to be more user friednly and error proof
func mysqlConnection() string{
	//Load file
	file, err := os.Open("conf/mysql.conf")
	if(err != nil){
		//Return default if no file was found
		fmt.Println("MySQL file not found")
		return "webapp:password@/webapp?charset=utf8"
	}
	
	//Creates new scanner
	s := bufio.NewScanner(file)
	
	//Reads lines and put it in vars
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
	
	//Return string compiled from vars
	return user+":"+pass+"@tcp("+addr+")/"+tabl+"?"+args
}

//prints error to console if error exists
func check(e error){
	if(e != nil){
		fmt.Println(e.Error())
	}
}

//Adds user to userinfo
func addUser(name string, pass string) bool{
	if(name == "" || pass == ""){return false}
	
	//Open mysql connection with connection args
	db, err := sql.Open("mysql", mysqlConnection())
	check(err)
	
	//Defer closing until function is returned, (GO best practice)
	defer db.Close()
	
	//Query db
	_, err = db.Query("INSERT INTO userinfo (username, password) VALUES ('"+name+"', '"+pass+"')")
	if(err != nil){
		fmt.Println(err.Error())
		return false
	}else{
		fmt.Println("User added "+name)
		return true
	}
}

//Creates new session if no exists
func newSession(session string) bool{
	if(session == ""){return false}
	
	//Open mysql connection with connection args
	db, err := sql.Open("mysql", mysqlConnection())
	check(err)
	
	//Defer closing until function is returned, (GO best practice)
	defer db.Close()
	
	//Query db
	_, err = db.Query("INSERT INTO sessions (session) VALUES ('"+session+"')")
	check(err)
	if(err == nil){
		fmt.Println("New session established, id: "+session)
		return true
	}
	fmt.Println("Existing session established, id: "+session)
	return false
}

//Updates user-column of sessions
func setUserForSession(session string, username string){
	if(session == ""){return}
	
	//Open mysql connection with connection args
	db, err := sql.Open("mysql", mysqlConnection())
	check(err)
	
	//Defer closing until function is returned, (GO best practice)
	defer db.Close()
	
	//Query db
	_, err = db.Query("UPDATE sessions SET user='"+username+"' WHERE session='"+session+"'")
	check(err)
}

func checkUserCredentials(username string, password string) bool{
	if(username == "" || password == ""){return false}
	
	//declares vars
	var user, pass string
	
	//Open mysql connection with connection args
	db, err := sql.Open("mysql", mysqlConnection())
	check(err)
	
	//Defer closing until function is returned, (GO best practice)
	defer db.Close()
	
	//Query db
	rows, err := db.Query("SELECT * FROM userinfo WHERE (username='"+username+"' AND password='"+password+"')")
	check(err)
	
	//Reads rows if not nil
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

//Gets username from sessions-table
func getUsernameBySessionId(session string) (string, bool){
	if(session == ""){return "", false}
	
	//Declare vars
	var username string
	
	//Open mysql connection with connection args
	db, err := sql.Open("mysql", mysqlConnection())
	check(err)
	
	//Defer closing until function is returned, (GO best practice)
	defer db.Close()
	
	//Query db
	rows, err := db.Query("SELECT user FROM sessions WHERE session='"+session+"'")
	check(err)
	
	//Reads rows if not nil
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

//Gets user-object from userinfo
func getUserByUsername(username string) (User, bool){
	if(username == ""){return User{name: "", pass: ""}, false}
	
	var user, pass string
	
	//Open mysql connection with connection args
	db, err := sql.Open("mysql", mysqlConnection())
	check(err)
	
	//Defer closing until function is returned, (GO best practice)
	defer db.Close()
	
	//Query db
	rows, err := db.Query("SELECT * FROM userinfo WHERE username='"+username+"'")
	check(err)
	
	//Reads rows if not nil
	if(rows != nil){
		for rows.Next(){
			err = rows.Scan(&user, &pass)
			if(user == username){
				//Initiates and returns new User-struct with username and password
				return User{name: user, pass: pass}, true
			}
		}
	}
	//Returnf empty User and false
	return User{name: "", pass: ""}, false
}
