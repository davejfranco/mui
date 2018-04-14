package store

import (
	"database/sql"
	"log"

	"github.com/davejfranco/mui/server/model"
)

//ServerStore struct
type ServerStore struct {
	Conn *sql.DB
}

//exist check if server is registered by ec2id
func (ss *ServerStore) exist(ec2Id string) bool {
	query := "SELECT ec2_id FROM server WHERE ec2_id=?"
	var ec2 string

	err := ss.Conn.QueryRow(query, ec2Id).Scan(&ec2)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal(err)
	}
	return true
}

//List all server in server table
func (ss *ServerStore) List() []model.Server {
	query := "SELECT * FROM server"

	row, err := ss.Conn.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	var servers []model.Server
	for row.Next() {
		var (
			ec2, ip string
		)

		err = row.Scan(&ec2, &ip)
		servers = append(servers, model.Server{Ec2Id: ec2, ServerIP: ip})

	}

	return servers
}

//New creates a new entry in server table
func (ss *ServerStore) New(server model.Server) error {

	query := "INSERT INTO `server` (ec2_id, server_ip, created) VALUES (?, ?, ?);"
	ss.Conn.Prepare(query)
	_, err := ss.Conn.Exec(server.Ec2Id, server.ServerIP)
	if err != nil {
		return err
	}
	return nil
}

//Delete server
func (ss *ServerStore) Delete(ec2Id string) error {
	delquery := "DELETE FROM server WHERE ec2_id=?"
	_, err := ss.Conn.Exec(delquery, ec2Id)
	if err != nil {
		return err
	}
	return nil
}
