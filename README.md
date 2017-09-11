# graphql_usr_srv
trying out go + graphql

* At least a golang mock-like graphql server for frontend testing
* Docker image from scratch with the compiled exe


working graphql queries:


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



## create user (doesn't actually create anything persistent)
mutation{
  createUser(username: "123", password: "123"){
    id
    username
  }
}
