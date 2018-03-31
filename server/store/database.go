package store

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//MysqlConfig struct takes all the parameters required to create connection to the mysql database
type MysqlConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DB       string
}

//URI method returns a string formatted uri like id:password@tcp(your-amazonaws-uri.com:3306)/dbname
//This method will latter be use to connect to a mysql server
func (config MysqlConfig) URI() string {
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

//Conn method returns an sql.DB object of the mysql connection
func (config MysqlConfig) Conn() *sql.DB {
	conn, err := sql.Open("mysql", config.URI())
	if err != nil {
		panic(err)
	}
	return conn
}

//Close method to terminate any given connection
func Close(conn *sql.DB) {
	conn.Close()
}
