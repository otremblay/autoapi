# AUTOAPI

Automatic api generation from an SQL database, complete with http API endpoint scaffolding code and preconditions checking.

---
## Prerequisites
* go >= 1.3 (see http://golang.org/)
* a mysql database
* go package [gorrilla/mux](http://www.gorillatoolkit.org/pkg/mux); (Dependency to be removed by #12)

```bash
go get github.com/gorilla/mux
```

## Installation
```bash
# Get package
go get is-a-dev.com/autoapi

# See that it works (if your go env is setup properly)
autoapi --help
```

---
## Create a new API
```bash
cd $GOPATH/src
mkdir autoapiapp && cd autoapiapp
autoapi -d="DB_NAME" -u="root" -h="localhost" -P="3306"
```
Once this step is complete your directory should now look like:
```
$GOPATH/src/whatever/
        bin/
            main.go (Main application binary, run this to start it)
        db/
            /mysql
                /DB_TABLES (Database queries themselves for mysql)
        dbi/
            /DB_TABLES
        http/
            /DB_TABLES
```

## Starting the API
From your project root:
```bash
go run bin/main.go
```
Your api will now be runnning on: http://localhost:8080 (by default)

## Additional Configuration

* How to change the port and host
* Other amazing features

## Working with your new project

##### # TODO need to write more guidelines and tips for how to use this for awesome things

* It is a good idea to use the routes generated as the base for your project, and include them in new packages instead of modifying them directly.

---
## Project Roadmap

* Finish constriant issues
* Complete Swagger
* Complete JSON-LD
* Onwards to glory

##### # TODO Need to add some more items here, not sure on priority

## Contributing
##### # TODO Need to write some guidelines you want for contributing to this project

## Contributors
* [Olivier Tremblay](https://git.is-a-dev.com/otremblay)
* [Colin Gagnon](ttps://github.com/colingagnon)