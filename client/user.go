package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const sudofile string = "/etc/sudoers.d/muisudo"

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
			if CheckGroup(v) {
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
	isavail := CheckUser(usr.Name)
	if !isavail {
		fmt.Println("user doesn't exist")
	} else {
		err := Execute("userdel", "-r", usr.Name)
		if err != nil {
			return err
		}
	}

	//Remove sudoer user
	err := DelSudoer(usr.Name)
	if err != nil {
		return err
	}
	return nil
}

func (usr User) Makesudo() error {
	file, err := os.OpenFile(sudofile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	sudostring := fmt.Sprintf("%s ALL=(ALL) NOPASSWD:ALL", usr.Name)
	reader := bufio.NewScanner(file)
	for reader.Scan() {
		if reader.Text() == sudostring {
			return nil
		}
	}

	if _, err = file.WriteString(fmt.Sprintf("%s\n", sudostring)); err != nil {
		return err
	}

	return nil

}
