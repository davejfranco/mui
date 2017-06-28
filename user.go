package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//type user to describe a user on linux systems
type User struct {
	Name      string   `json:"name"`
	Sudoer    bool     `json:"sudoer"`
	Groups    []string `json:"groups"`
	Publickey []string `json:"public_key"`
}

//Add user
func (usr User) Add() error {

	//_, err := user.Lookup(usr.name)
	isavail := checkUser(usr.Name)
	if !isavail {
		//This means that the user doesn'exist
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

		//create ssh directory
		sshdir := fmt.Sprintf("/home/%s/.ssh", usr.Name)
		err := os.Mkdir(sshdir, 0644)
		if err != nil {
			return err
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

	sudofile := "Cmnd_Alias     COMMANDS = %s\n%s ALL=(ALL) NOPASSWD:COMMANDS"

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
