# training-notebook

Server-backend for a simple training-notebook application. I use this as a 
testing ground for learning new things such as software designs/architectures, 
go libraries/frameworks, RESTful practices, authentication, database management etc.
The different branches in this repo show different libraries and architectures that
could be implemented. The 'main' branch displays the most robust iteration of the
server backend. 

# Overview

This is an HTTP server written in go and uses a SQLite database for permanent 
storage. There is a `frontend/` package that contains a react project, but it
will likely be obsoleted, as all javascript should be :). The app exposes 
RESTful endpoints for 'users' and 'sets', whose (swagger) documentation can be
found on the `/docs` server endpoint and rendered in the browser. Or you can 
view the spec in `docs/swagger.yaml`.

## Project Layout

The entry point to the server is `main.go`. `main` builds the server using the 
provided configs and loaded environment variables. The accepted configs are defined
in `config.go`. The loaded configs are used to create the server instance in the 
`builder.go` and `server.go` in the `server/` package. `main` starts the server 
and awaits a kill signal before graceful shutdown.

Models (objects used in the 'business logic') are defined in `models`, and contain
struct tags for json serialization in HTTP requests/responses.

Database interfaces are defined for the different models in `data/<model>.go`, e.g.
`data/user.go`. SQL statements and Database constructors are defined at the top of 
these files and the different operations are detailed below. Also defined in this 
package are mocks for the database interfaces, permitting flexible unit testing in 
the resource layer.

Resources and their associated CRUD operations are defined in the `api` package.
The path to a particular handler (or 'controller') adheres to the form 
`api/<resource>/<CRUD operation>`, for example, `api/sets/create`. Global middleware
(e.g. a latency logger) are defined directly in `api/` while resource-specific 
middlewares (e.g. user authentication/verification) are defined in their respective
resource packages.

`scripts/` contains useful shell scripts for testing, server deployment, 
go linting/vetting, and swagger spec generation.

`configs/` houses example configuration files written in yaml.

## Go Frameworks/Libraries 

Branch main features an implementation of a web-backend using the popular HTTP 
handling framework 'Gin', and alternatives such as 'gorilla' can be found on 
other branches. 

Also used in this project:
* go-swagger for API documentation generation
* golang-jwt for creating/parsing json web tokens
* joho/godotenv for loading environment variables
* crypto for password checking/hashing
* go-playground/validator for validating struct fields in json serialization

# Setup

You need a SQLite driver in order to run this server locally. I recommend: 
```
go get github.com/mattn/go-sqlite3  
go install github.com/mattn/go-sqlite3
```
You may also need to install or reinstall gcc. E.g. on Ubuntu:
```
sudo apt install --reinstall build-essential
```

# Testing

### Unit Testing

For all API endpoints (but not all middlewares) are unit tested, making use of a 
(manually created :P) Mock Databases defined in `data/`. Database operations are 
also unit tested by creating test-database instances. 

You can use the helper script `scripts/test.sh` to run individual package or all 
unit tests.

Run all unit tests:
```
./scripts/test.sh -a
```

Run package tests:
```
./scripts/test.sh -p api/sets
```

### Integration Testing

Integration tests are defined in `test/client_test.go`, and represent common use 
cases for the API. Integration tests operate by spinning up a test server locally
and sending HTTP requests to it through the clients defined in each test case.

You can use the helper script `scripts/test.sh` to start this server and run the 
integration tests like so:

```
./scripts/test.sh -i
```

### Testing & CI

This branch leverages Github Actions to build and test each commit, including integration
tests. The build is defined in `.github/workflows/test.yml`. Check out the 'Actions' tab
in github to see previous builds.

### Manual Testing

The server can be manually deployed locally by using the helper script `scripts/dev.sh`, 
which takes in a single argument for a yaml config (see examples of a server config in `configs/`).

E.g.
```
./scripts/dev.sh configs/config.yaml
```

This script isn't very useful yet, so you could also manually build and run the binary:
```
go build
./training-notebook --config=configs/config.yaml
```
However, I anticipate that as this project grows in complexity, deploying a server for testing
will require more robust configuration and setup.




