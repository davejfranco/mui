package store

import (
	"database/sql"
	"errors"
	"log"

	"github.com/davejfranco/mui/server/model"
	"github.com/davejfranco/mui/util"
)

type GroupStore struct {
	conn *sql.DB
}

//Exist return a boolean if the group exist in the group table
func (gs *GroupStore) Exist(groupname string) bool {

	var group string
	err := gs.conn.QueryRow("SELECT group_id FROM group WHERE group_id=?", groupname).Scan(&group)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal(err)
	}
	return true
}

//Get return a type Group if found on group table
func (gs *GroupStore) Get(groupname string) (model.Group, error) {
	query := "SELECT group_id, created FROM users WHERE group_id=?"

	var (
		groupid string
		created string
	)
	err := gs.conn.QueryRow(query, groupname).Scan(&groupid, &created)
	if err != nil {
		return model.Group{}, err
	}

	return model.Group{
		GroupID: groupid,
		Created: created,
	}, nil
}

//New method add new group to group table
func (gs *GroupStore) New(g model.Group) error {

	if len(g.GroupID) == 0 {
		return errors.New("Groupname is required")
	}

	//If group already exist, do nothing
	if gs.Exist(g.GroupID) {
		return nil
	}

	query := "INSERT INTO group (group_id, created) VALUES (?, ?)"
	_, err := gs.conn.Exec(query, g.GroupID, util.CurrentTime())
	if err != nil {
		return err
	}
	return nil
}

//Delete removes a group from a group table
func (gs *GroupStore) Delete(groupname string) error {
	delquery := "DELETE FROM group WHERE group_name=?"
	_, err := gs.conn.Exec(delquery, groupname)
	if err != nil {
		return err
	}
	return nil
}
