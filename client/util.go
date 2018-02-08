package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

const logfile string = "/var/log/mui.log"

func Log(severity string, message string) {

	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	//Close once finished
	defer f.Close()

	logger := log.New(f, fmt.Sprintf("%s - ", severity), log.Ldate|log.Ltime)
	logger.Println(message)
}

//Execute commands
func Execute(comand string, args ...string) error {
	if err := exec.Command(comand, args...).Run(); err != nil {
		return err
	}
	return nil
}

//Check if user exist
func CheckUser(username string) bool {
	_, err := user.Lookup(username)
	if err != nil {
		return false
	}
	return true
}

//Check if group exist
func CheckGroup(groupname string) bool {
	_, err := user.LookupGroup(groupname)
	if err != nil {
		return false
	}
	return true
}

func DelSudoer(user string) error {
	//scan the sudoer file with the line matching the user
	file, err := ioutil.ReadFile(sudofile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(file), "\n")
	sudostring := fmt.Sprintf("%s ALL=(ALL) NOPASSWD:ALL", user)
	//a = append(a[:i], a[i+1:]...)
	for i, line := range lines {
		if strings.Contains(line, sudostring) {
			lines = append(lines[:i], lines[i+1:]...)
			break
		}
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(sudofile, []byte(output), 0644)
	if err != nil {
		return err
	}
	return nil
}

func CheckSudo(user string) bool {
	file, err := ioutil.ReadFile(sudofile)
	if err != nil {
		return false
	}

	sudostring := fmt.Sprintf("%s ALL=(ALL) NOPASSWD:ALL", user)
	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if strings.Contains(line, sudostring) {
			return true
		}
	}
	return false
}

//Update SSH keys
func UpdateSSHKeys(user string, keys []string) {

	if len(keys) > 0 {
		err := Execute("rm", fmt.Sprintf("/home/%s/.ssh/authorized_keys", user))
		if err != nil {
			fmt.Println(err)
		}

		u := User{
			Name:      user,
			Publickey: keys,
		}

		err = u.AddPublicKey()
		if err != nil {
			fmt.Println(err)
		}
	}
}

//UnixUser to add a non IAM user
func UnixUser(user User) {

	//Check if user exist
	if CheckUser(user.Name) {
		fmt.Printf("user %s already exist, nothing to do here\n", user.Name)
		os.Exit(0)
	}
	err := user.Add()
	if err != nil {
		fmt.Printf("Cannot creater user: %s, Error: %s", user.Name, err)
	}
}

//delete user
func DelUser(user User) {

	if !CheckUser(user.Name) {
		fmt.Printf("user %s doesn't exist, nothing to do here\n", user.Name)
		os.Exit(0)
	}

	err := user.Del()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}

// func CurrentUsers() []User {

// 	home := "/Users"
// 	var users []User

// 	out, err := exec.Command("ls", home).Output()
// 	if err != nil {
// 		panic(err)
// 	}

// 	output := strings.Split(string(out), "\n")
// 	for _, user := range output[:len(output)-1] {
// 		users = append(users, User{Name: user})
// 	}

// 	return users

//}
