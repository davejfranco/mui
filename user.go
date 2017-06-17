/*
This is a user management library on golang -- duuh
Things I want to have on this library

- User creation
- Find User
- Delete User
- Make user sudoer
-
*/
package user

import (
	"fmt"
	"io/ioutil"
	"os"
	osuser "os/user"
	"strings"
)

type User struct {
	Name      string   `json:"name"`
	Sudoer    bool     `json:"sudoer"`
	Groups    []string `json:"groups"`
	Publickey string   `json:"public_key"`
}

//Check if user exist
func checkUser(username string) bool {
	_, err := osuser.Lookup(username)
	if err != nil {
		return false
	}
	return true
}

//Check if group exist
func checkGroup(groupname string) bool {
	_, err := osuser.LookupGroup(groupname)
	if err != nil {
		return false
	}
	return true
}

func (usr User) Add() error {

	//_, err := user.Lookup(usr.name)
	isavail := checkUser(usr.Name)
	if !isavail {
		//This means that the user doesn'exist
		cmd := "useradd"

		if len(usr.Groups) > 0 {
			//If the length of the field group is greater than zero check groups and then
			validGroup := make([]string, len(usr.Groups))
			for _, v := range usr.Groups {
				if checkGroup(v) {
					validGroup = append(validGroup, v)
				}
			}

			//Run useradd plus the groups that the user should be member
			err := Execute(cmd, "-G", strings.Join(validGroup, ","), "-m", usr.Name)
			if err != nil {
				return err //This should go to a log
			}

		}
		err := Execute(cmd, "-m", usr.Name)
		if err != nil {
			return err
		}

		//create ssh directory
		sshdir := fmt.Sprintf("/home/%s/.ssh", usr.Name)
		err := os.Mkdir(sshdir, 0644)
		if err != nil {
			//fmt.Println("Unable to add ssh directory")
			//os.Exit(1)
			return err
		}

		//Add SSH key, should this be a separate method?
		aerr := ioutil.WriteFile(sshdir+"/authorized_keys", []byte(usr.Publickey), 0600)
		if aerr != nil {
			fmt.Println(aerr)
		}

		//Change ownership to the user
		err := Execute("chown", "-R", string(fmt.Sprintf("%s:%s", usr.Name, usr.Name)), sshdir)
		//fmt.Println(string(fmt.Sprintf("%s:%s", usr.name)), sshdir)

		fmt.Println("Successfully user created")
	} else {
		fmt.Println("user already exist")
	}
}

func (usr User) Del() {
	//Check if user exist
	isavail := checkUser(usr.Name)
	if !isavail {
		fmt.Println("user doesn't exist")
	} else {
		Execute("userdel", "-r", usr.Name)
	}
}

//Create a sudo file for user
func (usr User) Makesudo() {

	//file := []byte("user ALL = NOPASSWD: ALL\nuser ALL=(ALL) NOPASSWD:ALL\n")
	file := fmt.Sprintf("%s ALL = NOPASSWD: ALL\n%s ALL=(ALL) NOPASSWD:ALL\n", usr.Name, usr.Name)

	err := ioutil.WriteFile(fmt.Sprintf("/etc/sudoers.d/%s", usr.Name), []byte(file), 0644)
	if err != nil {
		panic(err)
	}
}
