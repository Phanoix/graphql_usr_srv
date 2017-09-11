# graphql_usr_srv
trying out go + graphql

* At least a golang mock-like graphql server for frontend testing
* Docker image from scratch with the compiled exe


working graphql queries (listens on port 8000 by default):


## query test user
query{
  user(id: "1"){
    id
    username
  }
}

or just

query{
  user{
    id
    username
  }
}

## query test session
query{
  session(id: "1as6d546310asdf64@#9"){
    id
    expires
  }
}

## create user (doesn't actually create anything persistent)
mutation{
  createUser(username: "123", password: "123"){
    id
    username
  }
}

## login, get session
mutation{
  createSession(username: "test", password: "password"){
    id
    expires
  }
}
