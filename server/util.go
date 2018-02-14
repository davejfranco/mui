package main

import (
	"database/sql"
	"time"
)

func CurrentTime() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
