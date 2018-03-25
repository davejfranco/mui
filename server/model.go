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

//MysqlConfig struct takes all the parameters required to create connection to the mysql database
type MysqlConfig struct {
	user     string
	password string
	host     string
	port     int
	db       string
}

//URI method returns a string formatted uri like id:password@tcp(your-amazonaws-uri.com:3306)/dbname
//This method will latter be use to connect to a mysql server
func (config MysqlConfig) URI() string {
	//Mandatory database name
	if config.db == "" {
		fmt.Println("Database name must be specified")
		os.Exit(1)
	}

	var user, pass, address string

	if config.user != "" {
		user = config.user
	}

	user = config.user

	if config.password != "" {
		pass = fmt.Sprintf(":%s@", config.password)
	}

	if config.host != "" {
		if config.port != 0 {
			address = fmt.Sprintf("tcp(%s:%s)", config.host, strconv.Itoa(config.port))
		} else {
			address = fmt.Sprintf("tcp(%s)", config.host)
		}
	} else {
		address = ""
	}

	return strings.Join([]string{user, pass, address, "/", config.db}, "")
}

//Conn method returns an sql.DB object of the mysql connection
func (m MysqlConfig) Conn() *sql.DB {
	conn, err := sql.Open("mysql", m.URI())
	if err != nil {
		panic(err)
	}
	return conn
}

//User type
type User struct {
	Username  string `db:"username" json:"username"`
	Sudo      string `db:"sudo" json:"sudo"`
	Fullname  string `db:"full_name" json:"full_name"`
	Publickey string `db:"public_key" json:"public_key"`
	Created   string `db:"created" json:"created"`
}

//UserExist checks on database is the username exist in the users table
func UserExist(db *sql.DB, user string) bool {

	query := "SELECT username FROM users WHERE username=?"
	var username string

	err := db.QueryRow(query, user).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal(err)
	}
	return true
}

//GetUser returns a type User or error from the database based on a given username
func GetUser(db *sql.DB, user string) User {
	query := "SELECT * FROM users WHERE username=?"

	var (
		usr  string
		full sql.NullString //This field can be Null
		su   string
		pk   sql.NullString //This field can also be Null
		crea string
	)

	err := db.QueryRow(query, user).Scan(&usr,
		&full,
		&su,
		&pk,
		&crea)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}
		}
		log.Fatal(err)
	}

	usu := User{
		Username: user,
		Sudo:     su,
		Created:  crea,
	}

	if full.Valid {
		usu.Fullname = full.String
	}

	if pk.Valid {
		usu.Publickey = pk.String
	}

	return usu
}

//DelUser deletes user from users table
func DelUser(db *sql.DB, user string) error {
	delquery := "DELETE FROM users WHERE username=?"
	_, err := db.Exec(delquery, user)
	if err != nil {
		return err
	}
	return nil
}

//GetallUsers returns all users in users table
func GetallUsers(db *sql.DB) []User {

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}

	var allusers []User
	for rows.Next() {
		var (
			user string
			full sql.NullString //This field can be Null
			su   string
			crea string
			pk   sql.NullString //This field can also be Null
		)

		err = rows.Scan(&user, &full, &su, &crea, &pk)
		if err != nil {
			log.Fatal(err)
		}

		usr := User{
			Username: user,
			Sudo:     su,
			Created:  crea,
		}

		if full.Valid {
			usr.Fullname = full.String
		}

		if pk.Valid {
			usr.Publickey = pk.String
		}

		allusers = append(allusers, usr)
	}
	return allusers
}

//Group type
type Group struct {
	GroupID string `db:"group_id"`
	Created string `db:"created"`
}

//GroupExist return a boolean if the group exist in the group table
func GroupExist(db *sql.DB, groupname string) bool {

	var group string
	err := db.QueryRow("SELECT group_id FROM group WHERE group_id=?", groupname).Scan(&group)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal(err)
	}
	return true
}

//GetGroup return a type Group if found on group table
func GetGroup(db *sql.DB, groupname string) (Group, error) {
	query := "SELECT group_id, created FROM users WHERE group_id=?"

	var (
		groupid string
		created string
	)
	err := db.QueryRow(query, groupname).Scan(&groupid, &created)
	if err != nil {
		return Group{}, err
	}

	return Group{
		GroupID: groupid,
		Created: created,
	}, nil
}

//DelGroup function deletes a group from a group table
func DelGroup(db *sql.DB, groupname string) error {
	delquery := "DELETE FROM group WHERE group_name=?"
	_, err := db.Exec(delquery, groupname)
	if err != nil {
		return err
	}
	return nil
}

//Server type
type Server struct {
	Ec2Id    string `db:"ec2_id" json:"ec2_id"`
	ServerIP string `db:"server_ip" json:"server_ip"`
}

//Model type
type Model struct {
	conn *sql.DB
	User
	Group
	Server
}

//NewUser adds new user on "user" table
func (m *Model) NewUser() error {

	//If user exist do not continue
	if UserExist(m.conn, m.Username) {
		return nil
	}

	switch {
	case len(m.Username) == 0:
		return errors.New("username field is required")
	case len(m.Sudo) == 0:
		return errors.New("Sudo field is required")
	}

	insertquery := "INSERT INTO users (username, full_name, sudo, public_key, created) VALUES (?, ?, ?, ?, ?)"
	_, err := m.conn.Exec(insertquery, m.Username, NewNullString(m.Fullname), m.Sudo, NewNullString(m.Publickey), CurrentTime())
	if err != nil {
		return err
	}
	return nil

}

//UpdateUser method updates a given column of a user in users table
func (m *Model) UpdateUser(field, value string) error {
	//"UPDATE mui.users SET full_name="Carlos D. Rodriguez" WHERE username='crodriguez'"
	if field != "username" && field != "full_name" && field != "sudo" && field != "public_key" {
		return errors.New("Invalid field")
	}

	queryrun := func(q, v string) error {
		_, err := m.conn.Exec(q, v, m.Username)
		if err != nil {
			return err
		}
		return nil
	}

	switch {
	case field == "username":
		query := fmt.Sprintf("Update users SET username=? WHERE username=?")
		m.Username = value
		return queryrun(query, value)
	case field == "full_name":
		query := fmt.Sprintf("Update users SET full_name=? WHERE username=?")
		m.Fullname = value
		return queryrun(query, value)
	case field == "sudo":
		query := fmt.Sprintf("Update users SET sudo=? WHERE username=?")
		m.Sudo = value
		return queryrun(query, value)
	case field == "public_key":
		query := fmt.Sprintf("Update users SET public_key=? WHERE username=?")
		m.Publickey = value
		return queryrun(query, value)
	}
	return nil

}

//NewGroup method add new group to group table
func (m *Model) NewGroup() error {
	if len(m.GroupID) == 0 {
		return errors.New("Groupname is required")
	}

	//If group already exist, do nothing
	if GroupExist(m.conn, m.GroupID) {
		return nil
	}

	query := "INSERT INTO group (group_id, created) VALUES (?, ?)"
	_, err := m.conn.Exec(query, m.GroupID, CurrentTime())
	if err != nil {
		return err
	}
	return nil
}

//NewServer method create a new entry in server table
func (m *Model) NewServer() error {

	query := "INSERT INTO `server` (ec2_id, server_ip, created) VALUES (?, ?, ?);"
	m.conn.Prepare(query)
	_, err := m.conn.Exec(m.Ec2Id, m.ServerIP)
	if err != nil {
		return err
	}
	return nil
}
