package main

import (
	"database/sql"
	"time"
)

//CurrentTime function returns the current time in string format when is execute it
func CurrentTime() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
}

//NewNullString returns a type sql.NullString from a given string
//This function is required when there is a chance of a Null column on a database query
//Found this useful function on https://stackoverflow.com/questions/40266633/golang-insert-null-into-sql-instead-of-empty-string
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
