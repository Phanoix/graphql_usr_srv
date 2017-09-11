package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"

	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
)

func main() {
	fmt.Printf("Server starting...\n")
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))

	http.Handle("/query", &relay.Handler{Schema: schema})

	log.Fatal(http.ListenAndServe(":8000", nil))
	fmt.Printf("Server running...\n")
}



func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, %s world!", r.URL.Path[1:])
}


var schema *graphql.Schema

func init() {
	schema = graphql.MustParseSchema(Schema, &Resolver{})
}


var Schema = `
	schema {
		query: Query
		mutation: Mutation
	}
	# The query type, represents all of the entry points into our object graph
	type Query {
		user(id: ID = "1"): User
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
	input UserInput {
		# a unique username
		username: String!
		# user's password
		password: String!
	}
`


type userInput struct {
	Username	string
	Password	string
}

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

func (r *userResolver) Username() string {
	return r.u.Username
}


func (r *Resolver) CreateUser(args *struct {
	Username string
	Password  string
}) *userResolver{
	passwd, err := bcrypt.GenerateFromPassword([]byte("salt and such"+args.Password), 10)
	if err != nil { 
		fmt.Printf("user create error while generating Password")
		return nil
	}
	Password := base64.StdEncoding.EncodeToString(passwd)
	var createdUser = user{
		ID:			"2",
		Username:	args.Username,
		Password:	Password,
	}
	return &userResolver{&createdUser}
}









var page = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.css" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.1.0/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				return fetch("/query", {
					method: "post",
					body: JSON.stringify(graphQLParams),
					credentials: "include",
				}).then(function (response) {
					return response.text();
				}).then(function (responseBody) {
					try {
						return JSON.parse(responseBody);
					} catch (error) {
						return responseBody;
					}
				});
			}
			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
				document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`)