package main

import "testing"

func TestTypesAvailable(t *testing.T){
	var check_types [1]string
	check_types[0] = "User"

	testUser := user{
		ID:			"1",
		username:	"testuser",
		password:	"correct horse battery staple",
	}

	if testUser.username != "testuser" {
		panic("failed user type testuser");
	}
}