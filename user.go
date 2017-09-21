/*
 * User service
 */

package main

import (
 	"fmt"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-redis/redis"
)

type user struct {
	ID			string
	Username	string
	Password	string
}

var testUser = user{
	ID:			"1",
	Username:	"testuser",
	Password:	"correct horse battery staple",
}

type session struct {
	ID			string
	UserID		string
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


func getUserByID( ID string ) user {
	fetchedUser, err := fetchUser( ID )

	if err != nil{
		return testUser
	}

	return fetchedUser
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
	saveUser( createdUser )
	return &createdUser
}

func login( username string, pass string ) *session {
	return &testSession
}



func saveUser( newuser user ) {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ID :=  newuser.Username

	err := client.Set(ID+":username", newuser.Username, 0).Err()
	if err != nil {
		panic(err)
	}
	err = client.Set(ID+":password", newuser.Password, 0).Err()
	if err != nil {
		panic(err)
	}
}

func fetchUser( ID string ) (user, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	username, err := client.Get(ID+":username").Result()
	if err == redis.Nil {
		return user{}, err
	} else if err != nil {
		panic(err)
	}

	// fetch other fields as they're added the same way

	fetchedUser := user{
		ID:			ID,
		Username:	username,
	}

	return fetchedUser, nil
}