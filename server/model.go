package main

import (
	"database/sql"
	"errors"
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
	Sudo      string `db:"sudo" json:"sudo"`
	FullName  string `db:"full_name" json:"full_name"`
	PublicKey string `db:"public_key" json:"public_key"`
	Created   string `db:"created" json:"created"`
}

type UserModel struct {
	*User
	Conn *sql.DB
}

//Check if user exist
func (u UserModel) Exist() bool {

	query := "SELECT username FROM users WHERE username=?"
	var username string

	err := u.Conn.QueryRow(query, u.Username).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal(err)
	}
	return true
}

//Get specific user
func (u UserModel) Get() (User, error) {
	query := "SELECT full_name, sudo,  public_key, created FROM users WHERE username=?"

	var (
		fullname  sql.NullString //This field can be Null
		sudo      string
		publickey sql.NullString //This field can also be Null
		created   string
	)

	err := u.Conn.QueryRow(query, u.Username).Scan(&fullname,
		&sudo,
		&publickey,
		&created)

	if err != nil {
		return User{}, err
	}
	user := User{
		Username: u.Username,
		Sudo:     sudo,
		Created:  created,
	}

	if fullname.Valid {
		user.FullName = fullname.String
	}

	if publickey.Valid {
		user.PublicKey = publickey.String
	}

	return user, nil
}

//Add new user
func (u *UserModel) New() error {

	if u.Exist() {
		return nil
	}

	switch {
	case len(u.Username) == 0:
		return errors.New("Username field is required")
	case len(u.Sudo) == 0:
		return errors.New("Sudo field is required")
	}

	insertquery := "INSERT INTO users (username, full_name, sudo, public_key, created) VALUES (?, ?, ?, ?, ?)"
	_, err := u.Conn.Exec(insertquery, u.Username, NewNullString(u.FullName), u.Sudo, NewNullString(u.PublicKey), CurrentTime())
	if err != nil {
		return err
	}
	return nil
}

//Update user
func (u *UserModel) Update(field, value string) error {
	//"UPDATE mui.users SET full_name="Carlos D. Rodriguez" WHERE username='crodriguez'"
	if field != "username" && field != "full_name" && field != "sudo" && field != "public_key" {
		return errors.New("Invalid field")
	}

	queryrun := func(q, v string) error {
		_, err := u.Conn.Exec(q, v, u.Username)
		if err != nil {
			return err
		}
		return nil
	}

	switch {
	case field == "username":
		query := fmt.Sprintf("Update users SET username=? WHERE username=?")
		u.Username = value
		return queryrun(query, value)
	case field == "full_name":
		query := fmt.Sprintf("Update users SET full_name=? WHERE username=?")
		u.FullName = value
		return queryrun(query, value)
	case field == "sudo":
		query := fmt.Sprintf("Update users SET sudo=? WHERE username=?")
		u.Sudo = value
		return queryrun(query, value)
	case field == "public_key":
		query := fmt.Sprintf("Update users SET public_key=? WHERE username=?")
		u.PublicKey = value
		return queryrun(query, value)
	}
	return nil
}

//Delete user
func (u *UserModel) Del() error {
	delquery := "DELETE FROM users WHERE username=?"
	_, err := u.Conn.Exec(delquery, u.Username)
	if err != nil {
		return err
	}
	return nil
}

//Get all users
func (u UserModel) Getall() []User {

	rows, err := u.Conn.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}

	var allusers []User
	for rows.Next() {
		var (
			username  string
			fullname  sql.NullString //This field can be Null
			created   string
			publickey sql.NullString //This field can also be Null
		)

		err = rows.Scan(&username, &fullname, &created, &publickey)
		if err != nil {
			log.Fatal(err)
		}

		user := User{
			Username: username,
			Created:  created,
		}

		if fullname.Valid {
			user.FullName = fullname.String
		}

		if publickey.Valid {
			user.PublicKey = publickey.String
		}

		allusers = append(allusers, user)
	}
	return allusers
}

type Group struct {
	GroupId string `db:"groupid"`
	Users   []User
}

type GroupModel struct {
	Group
	Conn *sql.DB
}

//Server type
type Server struct {
	Ec2Id    string `db:"ec2_id" json:"ec2_id"`
	ServerIP string `db:"server_ip" json:"server_ip"`
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
