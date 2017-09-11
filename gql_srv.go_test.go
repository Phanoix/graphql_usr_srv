package main

import "testing"

func TestTypesAvailable(t *testing.T){
	// test user type existence
	testUser := user{
		ID:			"1",
		Username:	"testuser",
		Password:	"correct horse battery staple",
	}

	if testUser.Username != "testuser" {
		panic("failed user type testuser");
	}

	// test session type existence
	testSession := session{
		ID:			"634326846daf1",
		UserID:		"1",
		Username:	"testuser",
		Expires:	15488632,
	}

	if testSession.ID != "634326846daf1" {
		panic("failed session type testuser");
	}
}