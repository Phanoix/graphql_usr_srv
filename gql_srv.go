package main

import (
	"fmt"
//	"log"
	"net/http"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"

	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
)

func main() {
	fmt.Printf("Server starting...\n")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8000", nil)
	fmt.Printf("Server running...\n")
}



func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, %s world!", r.URL.Path[1:])
}


var Schema = `
	schema {
		query: Query
		mutation: Mutation
	}
	# The query type, represents all of the entry points into our object graph
	type Query {
		user(id: ID!): Human
	}
	# The mutation type, represents all updates we can make to our data
	type Mutation {
		createUser(username: String!, password: String!): User
	}
	# A user
	interface User {
		# The ID of the user
		id: ID!
		# The username
		username: String!
	}
	# The input object sent for creating a new user
	input ReviewInput {
		# a unique username
		username: String!
		# user's password
		password: String!
	}
`


type userInput struct {
	username	string
	password	string
}

type user struct {
	ID			graphql.ID
	username	string
	password	string
}

var testUser = user{
	ID:			"1",
	username:	"testuser",
	password:	"correct horse battery staple",
}


type Resolver struct{}

func (r *Resolver) User(args struct{ ID graphql.ID }) *userResolver {


	return &userResolver{&testUser}
}

type userResolver struct {
	u *user
}

func (r *userResolver) ID() graphql.ID {
	return r.u.ID
}

func (r *userResolver) username() string {
	return r.u.username
}


func (r *Resolver) CreateUser(args *struct {
	username string
	password  string
}) *userResolver{
	passwd, err := bcrypt.GenerateFromPassword([]byte("salt and such"+args.password), 10)
	if err != nil { 
		fmt.Printf("user create error while generating password")
		return nil
	}
	password := base64.StdEncoding.EncodeToString(passwd)
	var createdUser = user{
		ID:			"2",
		username:	args.username,
		password:	password,
	}
	return &userResolver{&createdUser}
}

