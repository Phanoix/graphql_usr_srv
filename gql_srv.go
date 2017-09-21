package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
)

func main() {
	fmt.Printf("Server started!\n")
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))

	http.Handle("/query", &relay.Handler{Schema: schema})

	log.Fatal(http.ListenAndServe(":8000", nil))
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
		user(id: String = "1"): User
		session(id: String!): Session
	}
	# The mutation type, represents all updates we can make to our data
	type Mutation {
		# user registration
		createUser(username: String!, password: String!): User
		# user login
		createSession(username: String!, password: String!): Session
	}
	# A user
	interface User {
		# The ID of the user
		id: String!
		# The username
		username: String!
	}
	# A session
	interface Session {
		# The ID of the user
		id: String!
		# The username
		expires: Float!
	}
	# The input object sent for creating a new user
	input UserInput {
		# a unique username
		username: String!
		# user's password
		password: String!
	}# The input object sent for logging in a user
	input SessionInput {
		# a unique username
		username: String!
		# user's password
		password: String!
	}
`

// inputs
type userInput struct {
	Username	string
	Password	string
}

type sessionInput struct {
	Username	string
	Password	string
}


type Resolver struct{}

// User resolving
func (r *Resolver) User(args struct{ ID string }) *userResolver {
	usr := getUserByID(args.ID)
	return &userResolver{&usr}
}

type userResolver struct {
	u *user
}

func (r *userResolver) ID() string {
	return r.u.ID
}
func (r *userResolver) Username() string {
	return r.u.Username
}

func (r *Resolver) CreateUser(args *struct {
	Username string
	Password  string
}) *userResolver{
	return &userResolver{addUser( args.Username, args.Password )}
}


// Session resolving
func (r *Resolver) Session(args struct{ ID string }) *sessionResolver {
	return &sessionResolver{&testSession}
}

type sessionResolver struct {
	s *session
}

func (r *sessionResolver) ID() string {
	return r.s.ID
}
func (r *sessionResolver) UserID() string {
	return r.s.UserID
}
func (r *sessionResolver) Username() string {
	return r.s.Username
}
func (r *sessionResolver) Expires() float64 {
	return float64(r.s.Expires)
}

func (r *Resolver) CreateSession(args *struct {
	Username string
	Password  string
}) *sessionResolver{
	return &sessionResolver{login( args.Username, args.Password )}
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