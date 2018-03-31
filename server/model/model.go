package model

//User type
type User struct {
	Username  string
	Sudo      bool
	Fullname  string
	Publickey string
	Created   string
}

func (user *User) GrantSudo() {
	user.Sudo = true
}

//Group type
type Group struct {
	GroupID string `db:"group_id"`
	Created string `db:"created"`
}

//Groups and users
type Userpergroup struct {
	GroupID string
	Users   []User
}

//Server type
type Server struct {
	Ec2Id    string `db:"ec2_id" json:"ec2_id"`
	ServerIP string `db:"server_ip" json:"server_ip"`
}

type GroupsinServer struct {
	Ec2Id  string
	Groups []Group
}

type UserStorer interface {
	Get() (User, error)
	List() ([]User, error)
	New() error
	Update() error
	Delete() error
}
