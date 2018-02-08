package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DB       string
}

//return uri connection format id:password@tcp(your-amazonaws-uri.com:3306)/dbname
func (config MysqlConfig) Uri() string {
	//Mandatory database name
	if config.DB == "" {
		fmt.Println("Database name must be specified")
		os.Exit(1)
	}

	var user, pass, address string

	if config.User != "" {
		user = config.User
	}

	user = config.User

	if config.Password != "" {
		pass = fmt.Sprintf(":%s@", config.Password)
	}

	if config.Host != "" {
		if config.Port != 0 {
			address = fmt.Sprintf("tcp(%s:%s)", config.Host, strconv.Itoa(config.Port))
		} else {
			address = fmt.Sprintf("tcp(%s)", config.Host)
		}
	} else {
		address = ""
	}

	return strings.Join([]string{user, pass, address, "/", config.DB}, "")
}

//Create Connection to database
func (db MysqlConfig) Conn() *sql.DB {
	conn, err := sql.Open("mysql", db.Uri())
	if err != nil {
		panic(err)
	}
	return conn
}

//User type
type User struct {
	Username  string `db:"username" json:"username"`
	FullName  string `db:"full_name" json:"full_name"`
	Created   string `db:"created" json:"created"`
	PublicKey string `db:"public_key" json:"publickey"`
}

type Group struct {
	GroupId string `db:"groupid"`
	Users   []User
}

//Server type
type Server struct {
	Ec2Id    string `db:"ec2_id" json:"ec2_id"`
	ServerIP string `db:"server_ip" json:"server_ip"`
}

type Model struct {
	MysqlConfig
}

//IfExist user
func (user *User) ifExist(conn *sql.DB) bool {

	query := "SELECT username FROM mui.users WHERE username=?"
	var username string

	err := conn.QueryRow(query, user.Username).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal(err)
	}
	return true
}

//Add new user
func (user *User) newUser(conn *sql.DB) error {

	if user.ifExist(conn) {
		return nil
	}

	insertquery := "INSERT INTO mui.users (username, full_name, created) VALUES (?, ?, ?)"
	//not sure if to set the current time in here
	user.Created = CurrentTime()
	_, err := conn.Exec(insertquery, user.Username, user.FullName, user.Created)
	if err != nil {
		return err
	}
	return nil
}

//modUser
func (user *User) modUser(conn *sql.DB) error {
	//"UPDATE mui.users SET full_name="Carlos D. Rodriguez" WHERE username='crodriguez'"
	modquery := "UPDATE mui.users SET username=?, full_name=? WHERE username=?"
	_, err := conn.Exec(modquery, user.Username, user.FullName, user.Username)
	if err != nil {
		return err
	}
	return nil
}

//delUser hate this comments
func (user *User) delUser(conn *sql.DB) error {
	delquery := "DELETE FROM mui.users WHERE username=?"
	_, err := conn.Exec(delquery, user.Username)
	if err != nil {
		return err
	}
	return nil
}

func getAllUsers(conn *sql.DB) []User {

	query := "SELECT * FROM mui.users"
	rows, err := conn.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	var allusers []User
	for rows.Next() {
		var (
			username string
			fullname string
			created  string
		)

		err = rows.Scan(&username, &fullname, &created)
		if err != nil {
			log.Fatal(err)
		}
		user := User{
			Username: username,
			FullName: fullname,
			Created:  created,
		}
		allusers = append(allusers, user)
	}
	return allusers

}

//Add node
func (srv Server) newNode(conn *sql.DB) error {

	query := "INSERT INTO `mui.server` (ec2_id, server_ip, created) VALUES (?, ?, ?);"
	conn.Prepare(query)
	result, err := conn.Exec(srv.Ec2Id, srv.ServerIP)
	if err != nil {
		return err
	}

	fmt.Println(result)
	return nil
}
