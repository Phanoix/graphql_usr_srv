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
		createUser(username: String!, password: String!, email: String): User
		# user login
		createSession(username: String!, password: String!): Session
	}
	# A user
	interface User {
		# The ID of the user
		id: String!
		# The username
		username: String!
		# User's email
		email: String
		# Date user was created
		registered: String!
		# Date and time of last login
		lastlogin: String
		# Is the user account active?
		active: Boolean!
		# Is the user an admin?
		admin: Boolean!
		# URL to fetch the user's avatar from
		avatarurl: String
		# is the user part of The Organization?!? will probably become a regular organization field at some point
		organization: Boolean!
	}
	# A session
	interface Session {
		# The ID of the session
		id: String!
		# When the session was created, in unix time
		created: Float!
		# When the session expires, in unix time
		expires: Float!
	}
	# The input object sent for creating a new user
	input UserInput {
		# a unique username
		username: String!
		# user's password
		password: String!
		# user's email
		email: String
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
	Email		string
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
func (r *userResolver) Email() *string {
	if r.u.Email == "" {
		return nil
	}
	return &r.u.Email
}
func (r *userResolver) Registered() string {
	return r.u.Registered.String()
}
func (r *userResolver) Lastlogin() *string {
	if r.u.LastLogin.String() == "" {
		return nil
	}
	lLogin := r.u.LastLogin.String()
	return &lLogin
}
func (r *userResolver) Active() bool {
	return r.u.Active
}
func (r *userResolver) Admin() bool {
	return r.u.Admin
}
func (r *userResolver) Avatarurl() *string {
	if r.u.AvatarURL == "" {
		return nil
	}
	return &r.u.AvatarURL
}
func (r *userResolver) Organization() bool {
	return r.u.Organization
}

func (r *Resolver) CreateUser(args *struct {
	Username string
	Password string
	Email *string
}) *userResolver{
	return &userResolver{addUser( args.Username, args.Password, args.Email )}
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
func (r *sessionResolver) Created() float64 {
	return float64(r.s.Created)
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