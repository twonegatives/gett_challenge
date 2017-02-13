# Gett challange

Golang code challange was to perform a list of tasks:

1. Build a data model for drivers and metrics;
2. Import sample data into datastore using Golang;
3. Build an API for accessing drivers and metrics using Beego framework.

## Briefly
It is hosted on `aws`, so you can give it a try like:
```
http://ec2-52-26-223-176.us-west-2.compute.amazonaws.com:8080/v1/driver
```

Watch `routes` file for inspiration on what else can be done there.

## Getting Started

You should obviously have a database (dump is done from postgres) and Golang installed to run this application. 
First, create a new database, load dump from the former folder

```
psql -U postgres gett < task_1/schema.sql
```

After that you can seed created database (pay attention to `task2/database.yml` file and having the required go packages imported)

```
cd task_2
go get gopkg.in/yaml.v2
go get github.com/lib/pq
go get github.com/jinzhu/gorm
go run main.go database.go  drivers.go metrics.go
```

At last, you're all set to fire up a local server. Make sure you moved `beegett` application folder into your `$GOPATH/src` folder
(beego fires warns if it's not there).

```
go get github.com/astaxie/beego
go get github.com/beego/bee
bee run beegett
```

## Tests

API is covered with endpoint tests for both drivers and metrics.
Tests should be launched from `beegett/tests` folder.

```
go test -v 
```
