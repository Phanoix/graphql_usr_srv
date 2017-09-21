/*
 * User service
 */

package main

import (
 	"fmt"
	"encoding/base64"
	"time"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-redis/redis"
)

type user struct {
	ID			string
	Username	string
	Password	string
}

var testUser = user{
	ID:			"default example user",
	Username:	"testuser",
	Password:	"correct horse battery staple",
}

type session struct {
	ID			string
	UserID		string
	Username	string
	Created		int64
	Expires		int64
}

// test data
// TODO: replace with an actual db connection
var testSession = session{
	ID:			"example session - not logged in",
	UserID:		"1",
	Username:	"testuser",
	Expires:	1506381787,
}


func getUserByID( ID string ) user {
	fetchedUser, err := fetchUser( ID )

	ID1 := base64.StdEncoding.EncodeToString( []byte(ID + time.Now().String()) )
	testUser.ID = ID1
	if err != nil{
		return testUser
	}

	return fetchedUser
}

func getSessionByID( ID string ) session {
	fetchedSession, err := fetchSession( ID )

	if err != nil{
		return testSession
	}

	return fetchedSession
}


func addUser( username string, pass string ) *user {

	passwd, err := bcrypt.GenerateFromPassword([]byte("lots of salt"+pass), 10)
	if err != nil { 
		fmt.Printf("user create error while generating Password")
		return nil
	}
	password := base64.StdEncoding.EncodeToString(passwd)
	var createdUser = user{
		ID:			username,
		Username:	username,
		Password:	password,
	}
	saveUser( createdUser )
	return &createdUser
}

func login( username string, pass string ) *session {
	// check for correct username + password
	loginUser, err1 := fetchUserPassword( username )
	if err1 != nil {
		return &testSession
	}

	pwd, err1 := base64.StdEncoding.DecodeString( loginUser.Password )
	if err1 != nil {
		return &testSession
	}

	if err := bcrypt.CompareHashAndPassword(pwd, []byte("lots of salt"+pass)); err != nil {
		return &testSession
	}
	// create new session
	ID := base64.StdEncoding.EncodeToString( []byte(username + time.Now().String()) )
	created := time.Now().Unix()
	expires := time.Now().Add(time.Hour * 24 * 30).Unix()

	newSession := session{
		ID:			ID,
		UserID:		username,
		Username:	username,
		Created:	created,
		Expires:	expires,
	}

	return &newSession
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

func saveSession( newsession session ) {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ID := newsession.ID

	err := client.Set("session:"+ID+":userid", newsession.UserID, 0).Err()
	if err != nil {
		panic(err)
	}
	err = client.Set("session:"+ID+":username", newsession.Username, 0).Err()
	if err != nil {
		panic(err)
	}
	err = client.Set("session:"+ID+":created", newsession.Created, 0).Err()
	if err != nil {
		panic(err)
	}
	err = client.Set("session:"+ID+":expires", newsession.Expires, 0).Err()
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

func fetchUserPassword( ID string ) (user, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	password, err := client.Get(ID+":password").Result()
	if err == redis.Nil {
		return user{}, err
	} else if err != nil {
		panic(err)
	}

	// fetch other fields as they're added the same way

	fetchedUser := user{
		ID:			ID,
		Password:	password,
	}

	return fetchedUser, nil
}

func fetchSession( ID string ) (session, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	expires, err := client.Get("session:"+ID+":expires").Result()
	if err == redis.Nil {
		return session{}, err
	} else if err != nil {
		panic(err)
	}
	expiresInt, err := strconv.Atoi(expires)
	if err != nil {
		panic(err)
	}

	created, err := client.Get("session:"+ID+":created").Result()
	if err == redis.Nil {
		return session{}, err
	} else if err != nil {
		panic(err)
	}
	createdInt, err := strconv.Atoi(created)
	if err != nil {
		panic(err)
	}


	// fetch other fields as they're added the same way

	fetchedSession := session{
		ID:			ID,
		Created:	int64(createdInt),
		Expires:	int64(expiresInt),
	}

	return fetchedSession, nil
}