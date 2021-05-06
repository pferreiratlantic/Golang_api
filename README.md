# Golang_api
Exercise1

## Context of the Exercise

It was a challenge to choose what language to use on this exercise.
Since it was placed as preference the Golang and / or VueJs, i decided to give a shot to Golang.
I decided to develop the API with Golang, using Postgresql as database engine and for the frontend 
page i have done a simple html page with jquery / ajax support to handle http requests between client
and server.
I tried to have in consideration the performance level of the API, since it had a csv file with over 
one million entries to populate database, i have decided to work in some sort of bulk operation. 
The frontend i have invested minimum time, its has a simple page, on the first section is possible to load
a csv from URL or if it goes empty the Api will use the provided csv. It is possible to get all user from a specific
country and since it was a huge load of entries i decided to make it with offsets to be easier to read.
At the end of the page you can see a table with number of users per country.





## Run locally

- Start postgres
- Execute script from dir schemas:
	psql -U postgres postgres -f schema.sql

- Configure environment

``` bash
$ source env-sample
```

- Build and starting API:

```bash
$ export GO111MODULE=on
$ export GOFLAGS=-mod=vendor
$ /usr/local/go/bin/go mod download
$ /usr/local/go/bin/go mod vendor
$ /usr/local/go/bin/go build -o Exercise1Api
$ killall -9 Exercise1Api -v
$ ./Exercise1Api &
```

Server API is listening on localhost:10000

## Testing

```bash
$ go test -v
=== RUN   TestCreateUser
--- PASS: TestCreateUser (0.00s)
=== RUN   TestGetUser
--- PASS: TestGetUser (0.00s)
=== RUN   TestGetNonExistentUser
--- PASS: TestGetNonExistentUser (0.00s)
=== RUN   TestEmptyUserTable
--- PASS: TestEmptyUserTable (0.00s)
=== RUN   TestGetCountries
--- PASS: TestGetCountries (0.00s)
=== RUN   TestEmptyCountryTable
--- PASS: TestEmptyCountryTable (0.00s)
PASS
ok  	github.com/pferreiratlantic/Golang_api	0.015s

```

