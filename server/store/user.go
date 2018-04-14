package store

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/davejfranco/mui/server/model"
	"github.com/davejfranco/mui/util"
)

type UserStore struct {
	Conn *sql.DB
}

//exist checks on database is the username exist in the users table
func (us *UserStore) exist(user string) bool {

	query := "SELECT username FROM users WHERE username=?"
	var username string

	err := us.Conn.QueryRow(query, user).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal(err)
	}
	return true
}

//Get returns a type User or error from the database based on a given username
func (us *UserStore) Get(user string) model.User {
	query := "SELECT * FROM users WHERE username=?"

	var (
		usr  string
		full sql.NullString //This field can be Null
		su   string
		pk   sql.NullString //This field can also be Null
		crea string
	)

	err := us.Conn.QueryRow(query, user).Scan(&usr,
		&full,
		&su,
		&pk,
		&crea)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}
		}
		log.Fatal(err)
	}

	sudo, _ := strconv.ParseBool(su)

	usu := model.User{
		Username: user,
		Sudo:     sudo,
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

//List returns all users in users table
func (us *UserStore) List() []model.User {

	rows, err := us.Conn.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}

	var allusers []model.User
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

		sudo, _ := strconv.ParseBool(su)
		usr := model.User{
			Username: user,
			Sudo:     sudo,
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

//New adds new user on "user" table
func (us *UserStore) New(user model.User) error {

	//If user exist do not continue
	if us.exist(user.Username) {
		return nil
	}

	if len(user.Username) == 0 {
		return errors.New("username is required")
	}

	insertquery := "INSERT INTO users (username, full_name, sudo, public_key, created) VALUES (?, ?, ?, ?, ?)"
	_, err := us.Conn.Exec(insertquery,
		user.Username,
		util.NewNullString(user.Fullname),
		strconv.FormatBool(user.Sudo),
		util.NewNullString(user.Publickey),
		util.CurrentTime())

	if err != nil {
		return err
	}
	return nil
}

//Update method updates a given column of a user in users table
func (us *UserStore) Update(user model.User) error {

	current := us.Get(user.Username)
	if current == (model.User{}) {
		return errors.New("not user found")
	}

	query := "UPDATE users SET full_name = ?, sudo = ?, public_key = ? WHERE username = ?"
	_, err := us.Conn.Exec(query, user.Fullname, strconv.FormatBool(user.Sudo), user.Publickey, user.Username)
	if err != nil {
		return err
	}
	return nil
}

//Delete deletes user from users table
func (us *UserStore) Delete(user string) error {
	delquery := "DELETE FROM users WHERE username=?"
	_, err := us.Conn.Exec(delquery, user)
	if err != nil {
		return err
	}
	return nil
}
