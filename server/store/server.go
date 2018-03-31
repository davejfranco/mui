package store

import (
	"database/sql"

	"github.com/davejfranco/mui/server/model"
)

type ServerStore struct {
	conn *sql.DB
}

func (ss *ServerStore) Exist(ec2Id string) bool {

	return true
}

//New creates a new entry in server table
func (ss *ServerStore) New(server model.Server) error {

	query := "INSERT INTO `server` (ec2_id, server_ip, created) VALUES (?, ?, ?);"
	ss.conn.Prepare(query)
	_, err := ss.conn.Exec(server.Ec2Id, server.ServerIP)
	if err != nil {
		return err
	}
	return nil
}

//Delete server
func (ss *ServerStore) Delete(ec2Id string) error {
	return nil
}
