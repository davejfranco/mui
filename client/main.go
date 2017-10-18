package main

import "fmt"

func main() {

	u := User{
		Name:   "dave",
		Sudoer: true,
	}

	err := u.Add()
	if err != nil {
		Log("Error", "Unable to add user")
	}
	Log("INFO", fmt.Sprintf("user %s successfully added", u.Name))

	err = u.Makesudo()
	if err != nil {
		Log("Error", "Unable to grant sudo")
	}
}
