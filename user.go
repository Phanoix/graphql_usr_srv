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

func getUserByID() user {
	return testUser
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