package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

//Execute commands
func Execute(comand string, args ...string) error {
	if err := exec.Command(comand, args...).Run(); err != nil {
		return err
	}
	return nil
}

//Check if user exist
func checkUser(username string) bool {
	_, err := user.Lookup(username)
	if err != nil {
		return false
	}
	return true
}

//Check if group exist
func checkGroup(groupname string) bool {
	_, err := user.LookupGroup(groupname)
	if err != nil {
		return false
	}
	return true
}

//User type
type User struct {
	Name      string   `json:"name"`
	Sudoer    bool     `json:"sudoer"`
	Groups    []string `json:"groups"`
	Publickey []string `json:"public_key"`
}

//Add user
func (usr User) Add() error {

	cmd := "useradd"

	//When are more than one group to add
	if len(usr.Groups) > 0 {
		//If the length of the field group is greater than zero check groups and then

		var validGroup []string
		for _, v := range usr.Groups {
			//Is an existing group on the system?
			if checkGroup(v) {
				validGroup = append(validGroup, v)
			}
		}
		//Run useradd plus the groups that the user should be member
		err := Execute(cmd, "-G", strings.Join(validGroup, ","), "-m", usr.Name)
		if err != nil {
			return err
		}

	} else {
		//add user without groups
		err := Execute(cmd, "-m", usr.Name)
		if err != nil {
			return err
		}
	}

	//Add Publickey
	if len(usr.Publickey) > 0 {
		usr.AddPublicKey()
	}

	return nil
}

//AddPublicKey to the user
func (usr User) AddPublicKey() error {

	//create ssh directory
	sshdir := fmt.Sprintf("/home/%s/.ssh", usr.Name)
	if _, err := os.Stat(sshdir); os.IsNotExist(err) {
		err := os.Mkdir(sshdir, 0644)
		if err != nil {
			return err
		}
	}

	//Add public keys
	f, err := os.OpenFile(sshdir+"/authorized_keys", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	for _, v := range usr.Publickey {
		if _, err = f.WriteString(v + "\n"); err != nil {
			panic(err)
		}
	}

	//Change ownership to the user
	err = Execute("chown", "-R", string(fmt.Sprintf("%s:%s", usr.Name, usr.Name)), sshdir)
	if err != nil {
		return err
	}
	return nil
}

//Del user
func (usr User) Del() error {

	//Check if user exist
	isavail := checkUser(usr.Name)
	if !isavail {
		fmt.Println("user doesn't exist")
	} else {
		err := Execute("userdel", "-r", usr.Name)
		if err != nil {
			return err
		}
	}

	//Remove sudoer file if exist
	if _, err := os.Stat(fmt.Sprintf("/etc/sudoers.d/%s", usr.Name)); err == nil {
		err := Execute("rm", fmt.Sprintf("/etc/sudoers.d/%s", usr.Name))
		if err != nil {
			return err
		}
	}
	return nil
}

//Makesudo to grant sudo priveledges
func (usr User) Makesudo(commands ...string) error {

	sudofile := "Cmnd_Alias     COMMANDS = %s\n%s ALL=(ALL) NOPASSWD:COMMANDS\n"

	addsudoer := func(cmd string) error {
		file := fmt.Sprintf(sudofile, cmd, usr.Name)
		err := ioutil.WriteFile(fmt.Sprintf("/etc/sudoers.d/%s", usr.Name), []byte(file), 0644)
		if err != nil {
			return err
		}
		return nil
	}

	//if sudo should allow ALL commands
	if commands[0] == "ALL" || commands[0] == "all" {
		err := addsudoer(strings.ToUpper(commands[0]))
		return err
	}
	err := addsudoer(strings.Join(commands, ","))
	return err

}
