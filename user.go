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
	"strings"
)

type User struct {
	Name      string   `json:"name"`
	Sudoer    bool     `json:"sudoer"`
	Groups    []string `json:"groups"`
	Publickey []string `json:"public_key"`
}

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

/*
Me playing around
func main() {
	//Add user
	public_key := []string{
		`ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDNSloNNYSsk6ULLCbsPh5US47iq0dtYzE8gpUJeJrFDerGCrpMkfJdOtF78VS80lkl0M5p7GS0Sm3E/4tldXmyvo1GM/zJnBw7VuJ1Z7FI2x4hBu6z3eLBw8BWV0edrMKe2EXhHkbkly2fTtD7XKDr53coZ8G8/vtchEYtTnZ18FrAhE9SaWth20eQ47liJjK/sW0U3RmsL+5M2yxPt5LyftxxrAKg+YDyfbhijfvmvwrUVoyPWI3p9ndLhig4BciK5IUaUnDtZIPQmXxYLTdW5fgi+A8AeS66d3uiwWj3Au6yii1xwsQE/6YnzyudEVuBGtjRneTP8Yck81m4sxozZSp6++W0BhEQrWifzilcPLNr62myG2pUtidHlZtd9qEGA4sK3qbzsXF9ku0F0kLtjvdXJeRo6breHhQuvBLGZWDbZuqeM/tryp/75VyAOkI3fX0qTH6VW3Y881Z1ah9AqwcAIeGj1dyvuPGg2v7wxn+7BM5RwDLfqwAXUmIqoCzPvwxV2/uhgHTH4ZDoqyshzaqcKJhf+Py/DESAT3MQl+qnNzc5Uz4YPXj61dwmOUvkffPkzzrNEprm+vFn6AnWUBGWRHuK3k61wqW6BIJuNYn5nyKYJGLOUVVEZMulgSnI4FajOb8AoKs30EAITNMbC5n3T+DMNyOAAZDjQ+gdIQ== groupon@C02NLP3YG3QC.group.on`,
		`ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDNSloNNYSsk6ULLCbsPh5US47iq0dtYzE8gpUJeJrFDerGCrpMkfJdOtF78VS80lkl0M5p7GS0Sm3E/4tldXmyvo1GM/zJnBw7VuJ1Z7FI2x4hBu6z3eLBw8BWV0edrMKe2EXhHkbkly2fTtD7XKDr53coZ8G8/vtchEYtTnZ18FrAhE9SaWth20eQ47liJjK/sW0U3RmsL+5M2yxPt5LyftxxrAKg+YDyfbhijfvmvwrUVoyPWI3p9ndLhig4BciK5IUaUnDtZIPQmXxYLTdW5fgi+A8AeS66d3uiwWj3Au6yii1xwsQE/6YnzyudEVuBGtjRneTP8Yck81m4sxozZSp6++W0BhEQrWifzilcPLNr62myG2pUtidHlZtd9qEGA4sK3qbzsXF9ku0F0kLtjvdXJeRo6breHhQuvBLGZWDbZuqeM/tryp/75VyAOkI3fX0qTH6VW3Y881Z1ah9AqwcAIeGj1dyvuPGg2v7wxn+7BM5RwDLfqwAXUmIqoCzPvwxV2/uhgHTH4ZDoqyshzaqcKJhf+Py/DESAT3MQl+qnNzc5Uz4YPXj61dwmOUvkffPkzzrNEprm+vFn6AnWUBGWRHuK3k61wqW6BIJuNYn5nyKYJGLOUVVEZMulgSnI4FajOb8AoKs30EAITNMbC5n3T+DMNyOAAZDjQ+gdIQ== groupon@C02NLP3YG3QC.group.com`,
		`ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDNSloNNYSsk6ULLCbsPh5US47iq0dtYzE8gpUJeJrFDerGCrpMkfJdOtF78VS80lkl0M5p7GS0Sm3E/4tldXmyvo1GM/zJnBw7VuJ1Z7FI2x4hBu6z3eLBw8BWV0edrMKe2EXhHkbkly2fTtD7XKDr53coZ8G8/vtchEYtTnZ18FrAhE9SaWth20eQ47liJjK/sW0U3RmsL+5M2yxPt5LyftxxrAKg+YDyfbhijfvmvwrUVoyPWI3p9ndLhig4BciK5IUaUnDtZIPQmXxYLTdW5fgi+A8AeS66d3uiwWj3Au6yii1xwsQE/6YnzyudEVuBGtjRneTP8Yck81m4sxozZSp6++W0BhEQrWifzilcPLNr62myG2pUtidHlZtd9qEGA4sK3qbzsXF9ku0F0kLtjvdXJeRo6breHhQuvBLGZWDbZuqeM/tryp/75VyAOkI3fX0qTH6VW3Y881Z1ah9AqwcAIeGj1dyvuPGg2v7wxn+7BM5RwDLfqwAXUmIqoCzPvwxV2/uhgHTH4ZDoqyshzaqcKJhf+Py/DESAT3MQl+qnNzc5Uz4YPXj61dwmOUvkffPkzzrNEprm+vFn6AnWUBGWRHuK3k61wqW6BIJuNYn5nyKYJGLOUVVEZMulgSnI4FajOb8AoKs30EAITNMbC5n3T+DMNyOAAZDjQ+gdIQ== groupon@C02NLP3YG3QC.group.cl`,
	}
	user := User{
		Name:      "dfranco",
		Sudoer:    true,
		Groups:    []string{"adm", "operator"},
		Publickey: public_key,
	}

	//Add user
	err := user.Add()
	if err != nil {
		fmt.Printf("Error is %s", err)
	}

}
*/
