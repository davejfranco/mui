package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

//add user from IAM
func addIAMUser(user User) {

	//Check if user exist
	if checkUser(user.Name) {
		fmt.Printf("user %s already exist, nothing to do here", user.Name)
		os.Exit(0)
	}

	userKeys, err := IamUserKeys(user.Name)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	user.Publickey = userKeys

	err = user.Add()
	if err != nil {
		fmt.Printf("User cannot be created, Error: %s", err)
	}
}

func upateIamPublicKeys(user User) {

	if !checkUser(user.Name) {
		fmt.Printf("user %s doesn't exist, nothing to do here\n", user.Name)
		os.Exit(0)
	}

	iamUser, err := IamUserKeys(user.Name)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	authorizedfile := fmt.Sprintf("/home/%s/.ssh/authorized_keys", user.Name)
	if _, err := os.Stat(authorizedfile); os.IsExist(err) {
		Execute("rm", authorizedfile)
		fmt.Println("file deleted")
	}

	err = iamUser.AddPublicKey()
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

}

//add unix user
func addUnixUser(user User) {

	//Check if user exist
	if checkUser(user.Name) {
		fmt.Printf("user %s already exist, nothing to do here\n", user.Name)
		os.Exit(0)
	}
	err := user.Add()
	if err != nil {
		fmt.Printf("Cannot creater user: %s, Error: %s", user.Name, err)
	}
}

//grant sudo
func sudoCommands(user User, commands ...string) {

	err := user.Makesudo(strings.Join(commands, ","))
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}

//delete user
func delUser(user User) {

	if !checkUser(user.Name) {
		fmt.Printf("user %s doesn't exist, nothing to do here\n", user.Name)
		os.Exit(0)
	}

	err := user.Del()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

}

func myUsage() {
	fmt.Println("mui - Management User Interface")
	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", "mui")
	flag.PrintDefaults()
}

func main() {

	username := flag.String("user", "", "user reference. (Required)")
	aws := flag.Bool("aws", false, "Look existing user on AWS IAM")
	del := flag.Bool("del", false, "Remove user")
	sudo := flag.Bool("sudo", true, "Does the user should have sudo privileges")
	help := flag.Bool("help", false, "")
	//updatekey := flag.Bool("uk", false, "update IAM Public SSH keys")
	flag.Parse()

	user := User{
		Name:   *username,
		Sudoer: *sudo,
	}
	if *username == "" {
		myUsage()
		os.Exit(1)
	}

	if *help || len(os.Args) < 1 {
		myUsage()

	}

	switch {
	case *aws: // && !*updatekey:
		addIAMUser(user)
		if *sudo {
			user.Makesudo("all")
		}
	case *del:
		delUser(user)
	// case *aws && *updatekey:
	// 	upateIamPublicKeys(user)
	default:
		fmt.Println("Unrecognize combination of flags")
		myUsage()

	}

}
