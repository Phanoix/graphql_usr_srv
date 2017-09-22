# graphql_usr_srv
trying out go + graphql

* At least a golang mock-like graphql server for frontend testing - done! now with a simple redis user store
* Docker image from scratch with the compiled exe - done? not from scratch though


working graphql queries (listens on port 8000 by default):


## query test user
query{
  user(id: "1"){
    id
    username
    email
    registered
    lastlogin
    active
    admin
    avatarurl
    organization
  }
}

or just

query{
  user{
    id
    username
    email
    registered
    lastlogin
    active
    admin
    avatarurl
    organization
  }
}

## query test session
query{
  session(id: "1as6d546310asdf64@#9"){
    id
    created
    expires
  }
}

## create user
mutation{
  createUser(username: "123", password: "123", email: "a@b.c"){
    id
    username
    email
    registered
    lastlogin
    active
    admin
    avatarurl
    organization
  }
}

## login, get session
mutation{
  createSession(username: "test", password: "password"){
    id
    created
    expires
  }
}
