autoapi
=======

Automatic api generation from an SQL database, complete with http API endpoint scaffolding code and preconditions checking.


Installation
============

Because I'm a lazy bum, I have not gotten around to setting this up with a proper go-gettable url, so for now, you have to do custom stuff to install it.

1- clone repo to $GOPATH/src/is-a-dev.com/autoapi
2- inside that dir, run go get.
3- same dir, run go install


Usage
=====

autoapi <connectionstring> <dbname>

connectionstring: 'tcp:127.0.0.1:3306*databasename/username/password' <-- this is ugly as all hell and needs to go, I agree with you, now shush.
dbname: 		  your db name

A tiny guide:

cd in a new go project (empty dir under $GOPATH/src/wherever)
autoapi <connectionstring> <dbname>
if everything went according to plan, your dir now looks like this:

$GOPATH/src/whatever/
---bin/
------main.go
---http/
------several folders and/or depending on your db tables
---dbi/
------several folders and/or depending on your db tables
---db/
------mysql/
---------several folders and/or depending on your db tables.

You can just do 'go run bin/main.go <connectionstring> <dbname>' and you should have a shiny REST api boot up on port 8080!
