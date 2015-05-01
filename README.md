autoapi
=======

Automatic api generation from an SQL database, complete with http API endpoint scaffolding code and preconditions checking.


Installation
============

## Prerequisites
* go >= 1.3 (see http://golang.org/)
* a mysql database


## Procedure
    go get is-a-dev.com/autoapi

Usage
=====

    autoapi <dbhost> <dbname> <dbuser>

    dbhost:            ip or hostname
    dbname: 		   your db name
	dbuser:            database user that has access to read the information schema.

A tiny guide:

cd in a new go project (empty dir under $GOPATH/src/wherever)
autoapi <dbhost> <dbname> <dbuser>
if everything went according to plan, your dir now looks like this:

    $GOPATH/src/whatever/
    ├──bin/
    │  └──main.go
    ├──http/
    │  └──several folders and/or depending on your db tables
    ├──dbi/
    │  └──several folders and/or depending on your db tables
    └──db/
       └──mysql/
          └──several folders and/or depending on your db tables.

You can just do 'go run bin/main.go <dbhost> <dbname> <dbuser>' and you should have a shiny REST api boot up on port 8080!
