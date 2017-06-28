package main

import (
	"os/exec"
	"os/user"
)

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

//Execute commands
func Execute(comand string, args ...string) error {
	if err := exec.Command(comand, args...).Run(); err != nil {
		return err
	}
	return nil
}
