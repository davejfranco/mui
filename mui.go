package main

import (
	"flag"
	"fmt"
	"os"
)

func myUsage() {
	fmt.Println("mui - Management User Interface")
	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", "mui")
	flag.PrintDefaults()
}

func main() {

	username := flag.String("user", "", "user reference. (Required)")
	aws := flag.Bool("aws", false, "Look existing user on AWS IAM")
	sudo := flag.Bool("sudo", true, "Does the user should have sudo privileges")
	help := flag.Bool("help", false, "")
	flag.Parse()

	if *username == "" {
		myUsage()
		os.Exit(1)
	}

	if *help || len(os.Args) < 1 {
		myUsage()

	}

	u := User{
		Name:   *username,
		Sudoer: *sudo,
	}

}

// func main() {
//
// 	me := "dfranco"
// 	user, err := iamUser(me)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	//become sudo
// 	user.Sudoer = true
//
// 	user.Add()
// }
