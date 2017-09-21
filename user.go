/*
 * User service
 */

package main

import (
 	"fmt"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"

	"github.com/neelance/graphql-go"
)

type user struct {
	ID			graphql.ID
	Username	string
	Password	string
}

var testUser = user{
	ID:			"1",
	Username:	"testuser",
	Password:	"correct horse battery staple",
}

type session struct {
	ID			graphql.ID
	UserID		graphql.ID
	Username	string
	Expires		int
}

// test data
// TODO: replace with an actual db connection
var testSession = session{
	ID:			"1as6d546310asdf64@#9",
	UserID:		"1",
	Username:	"testuser",
	Expires:	1506381787,
}


func getUserByID() user {
	return testUser
}

func getSessionByID() session {
	return testSession
}


func addUser( username string, pass string ) *user {
	passwd, err := bcrypt.GenerateFromPassword([]byte("lots of salt"+pass), 10)
	if err != nil { 
		fmt.Printf("user create error while generating Password")
		return nil
	}
	password := base64.StdEncoding.EncodeToString(passwd)
	var createdUser = user{
		ID:			"2",
		Username:	username,
		Password:	password,
	}
	return &createdUser
}

func login( username string, pass string ) *session {
	return &testSession
}