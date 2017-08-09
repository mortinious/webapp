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

var db = initMySqlConnection()

var prepQ = make(map[string]*sql.Stmt)

//Declares User struct
type User struct{
	name string
	pass string
}

//Get mysql connection args from mysql.conf
//TODO: Update config file to be more user friednly and error proof
func initMySqlConnection() (*sql.DB){
	
	connection := "webapp:password@/webapp?charset=utf8"
	
	//Load file
	file, err := os.Open("conf/mysql.conf")
	if(err != nil){
		//Return default if no file was found
		fmt.Println("MySQL file not found")
		return nil
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
	connection = (user+":"+pass+"@tcp("+addr+")/"+tabl+"?"+args)
	db, err := sql.Open("mysql", connection)
	if(err != nil){
		fmt.Println(err.Error())
		return nil
	}
	
	return db
}

func prepareQueries(){
	//Prepare Queries
	prepareQ("addUser", "INSERT INTO userinfo (username, password) VALUES (?, ?)")
	prepareQ("newSession", "INSERT INTO sessions (session) VALUES (?)")
	prepareQ("setUserForSession", "UPDATE sessions SET user=? WHERE session=?")
	prepareQ("checkUserCred", "SELECT * FROM userinfo WHERE (username=? AND password=?)")
	prepareQ("getUserBySession", "SELECT user FROM sessions WHERE session=?")
	prepareQ("checkUsername", "SELECT 1 FROM userinfo WHERE username=?")
}

func prepareQ(e string, q string){
	//Prepare statement and add it to map
	stmt, err := db.Prepare(q)
	if(err == nil){
		prepQ[e] = stmt
	}
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
	
	db.Begin()
	
	//Query db
	_, err := prepQ["addUser"].Exec(name, pass)
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
	
	//Query db
	_, err := prepQ["newSession"].Exec(session)
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
	
	//Query db
	_, err := prepQ["setUserForSession"].Exec(username, session)
	fmt.Println("Setting user "+username+" for session "+session)
	check(err)
}

func checkUserCredentials(username string, password string) bool{
	if(username == "" || password == ""){return false}
	
	//declares vars
	var user, pass string
	
	//Query db
	rows, err := prepQ["checkUserCred"].Query(username, password)
	fmt.Println("Checking credentials for "+username+" with password "+password)
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
	
	//Query db
	row := prepQ["getUserBySession"].QueryRow(session)
	
	err := row.Scan(&username)
	if(err == nil){
		fmt.Println("-- Got: "+username)
		return username, true
	}
	return "", false
}

//Gets username from sessions-table
func checkUsername(username string) bool{
	if(username == ""){return false}
	
	var exists bool
	
	//Query db
	err := prepQ["checkUsername"].QueryRow(username).Scan(&exists)

	return err != sql.ErrNoRows
}

