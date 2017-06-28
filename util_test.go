package main

import (
	"testing"
)

func TestValidUser(t *testing.T) {
	want := true
	got := checkUser("root")
	if got != want {
		t.Fail()
	}
}

func TestUnvalidUser(t *testing.T) {
	want := false
	got := checkUser("dfranco")
	if got != want {
		t.Fail()
	}
}

func TestValidGroup(t *testing.T) {
	want := true
	got := checkGroup("admin")
	if got != want {
		t.Fail()
	}
}

func TestNotValidGroup(t *testing.T) {
	want := false
	got := checkGroup("root")
	if got != want {
		t.Errorf("got %T but want %T", got, want)
	}
}

func Test_nilexec(t *testing.T) {
	got := Execute("ls")
	if got != nil {
		t.Fail()
	}
}

func Test_failexec(t *testing.T) {
	got := Execute("ls -a")
	if got == nil {
		t.Fail()
	}

}
