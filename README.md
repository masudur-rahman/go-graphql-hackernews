# go-graphql-hackernews

## Project Setup

### Initialize go modules file:
```shell
go mod init github.com/masudur-rahman/hackernews
```

### Initialize a gqlgen project
```shell
go get -d github.com/99designs/gqlgen
go run github.com/99designs/gqlgen init
```

### Generate Go codes for the graphql schemas
```shell
go run github.com/99designs/gqlgen generate
```

#### Sometimes `gqlgen generate` doesn't work properly. You may need to run the followings
```shell
 printf '// +build tools\npackage tools\nimport (_ "github.com/99designs/gqlgen"\n _ "github.com/99designs/gqlgen/graphql/introspection")' | gofmt > tools.go\n\ngo mod tidy
 rm -rf vendor
 go run github.com/99designs/gqlgen generate
 go mod tidy && go mod vendor
```

**Note**: If you are getting `validation failed: packages.Load` error. It may occur, because `gqlgen` uses todo project as starter template.
To get rid of this error, edit `graph/schema.resolvers.go` file and delete functions `CreateTodo` and `Todos`. Now run the command again.


### Run MySQL Server from docker
```shell
docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=mysql -e MYSQL_DATABASE=hackernews -d mysql:latest
```
It will start a `MySQL Server` with a `root:mysql` user and a database named `hackernews`

### To terminate the mysql docker, just run the following command
```shell
docker rm mysql -f
```

### Exec into database
```shell
docker exec -it mysql bash

> mysql -u root -p
>> mysql # password

# Basic pre commands
> show databases;
> use hackernews;
> show tables;
> describe Users;
```

### Models and Migrations
We need to create migrations for our app so every time our app runs it creates tables it needs to work properly, we are going to use `golang-migrate` package.
Create a folder structure for our database files in the project root directory:
```
go-graphql-hackernews
--internal
----pkg
------db
--------migrations
----------mysql
```
##### Install go mysql driver and golang-migrate packages then create migrations:
```shell
go get -u github.com/go-sql-driver/mysql
go build -tags 'mysql' -ldflags="-X main.Version=1.0.0" -o $GOPATH/bin/migrate github.com/golang-migrate/migrate/v4/cmd/migrate/
cd internal/pkg/db/migrations/
migrate create -ext sql -dir mysql -seq create_users_table
migrate create -ext sql -dir mysql -seq create_links_table
```

`migrate` command will create two files for each migration ending with .up and .down; up is responsible for applying migration and down is responsible for reversing it.
Open `000001_create_users_table.up.sql` and add table for our users:
```sql
CREATE TABLE IF NOT EXISTS Users(
    ID INT NOT NULL UNIQUE AUTO_INCREMENT,
    Username VARCHAR (127) NOT NULL UNIQUE,
    Password VARCHAR (127) NOT NULL,
    PRIMARY KEY (ID)
)
```

in `000002_create_links_table.up.sql`:
```sql
CREATE TABLE IF NOT EXISTS Links(
    ID INT NOT NULL UNIQUE AUTO_INCREMENT,
    Title VARCHAR (255) ,
    Address VARCHAR (255) ,
    UserID INT ,
    FOREIGN KEY (UserID) REFERENCES Users(ID) ,
    PRIMARY KEY (ID)
)
```
We need one table for saving links and one table for saving users. 

#### If you want to migrate our tables from the CLI, then run the following `migrate` command from the project root directory
```shell
migrate -database mysql://root:dbpass@/hackernews -path internal/pkg/db/migrations/mysql up
```
**NB**: The migration is added to the starting of the `GraphQL` Server. So, we don't need to run it manually.


## Start the `GraphQL` Server
```shell
go run server/server.go
```

### Now open http://localhost:8080 to `Experience` the `hackernews` server
You will find a `GraphQL` Playground here. You can test the defined `queries` and `mutations` here.
Here's some example `queries` and `mutations`.

```graphql
mutation createUser {
  createUser(input: {username: "masud", password: "password"})
}
```
```graphql
mutation login {
  login(input: {username: "masud", password: "password"})
}
```
```graphql
mutation createLink {
  createLink(input: {title: "something", address: "somewhere"}) {
    title
    address
    id
    user {
      name
    }
  }
}
```
- With a Additional Request Header
    ```json
    {
        "Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Njg3ODgyODAsInVzZXJuYW1lIjoibWFzdWQifQ.3G0dggzOSXimbszXsD1yrTGlXgANhGyWrbf2S9vPHRU"
    }
    ```
```graphql
query links {
  links {
    title
    address,
    id
    user {
      name
    }
  }
}
```
