package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	url := MysqlConfig{
		User:     "muiadmin",
		Password: "muipass",
		Host:     "localhost",
		DB:       "mui",
	}


	conn := DBConn(url.Uri())
	defer conn.Close()

	me := User{
		Username: "dfranco",
		FullName: "Dave Franco",
	}

	//Create me
	fmt.Println(me.NewUser(conn))

}

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
func DBConn(uri string) *sql.DB {

	db, err := sql.Open("mysql", uri)
	if err != nil {
		panic(err)
	}

	return db
}

type User struct {
	Username string `db:"username"`
	FullName string `db:"full_name"`
}

//Add new user
func (usr User) NewUser(conn *sql.DB) error {

	insertquery := "INSERT INTO mui.users (username, full_name) VALUES (?, ?)"
	_, err := conn.Exec(insertquery, usr.Username, usr.FullName)
	if err != nil {
		return err
	}
	return nil
}

func 

type Group struct {
	Name string `db:"name"`
}

type Server struct {
	Ec2Id    string `db:"ec2_id"`
	ServerIp string `db:"server_ip"`
	UserId   int    `db:user_id`
}

func (srv Server) Add(conn *sql.DB) error {

	query := "INSERT INTO `server` (ec2_id, server_ip) VALUES (?, ?);"
	conn.Prepare(query)
	result, err := conn.Exec(srv.Ec2Id, srv.ServerIp)
	if err != nil {
		return err
	}

	fmt.Println(result)
	return nil

}
