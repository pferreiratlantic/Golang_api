# Golang_api
Exercise1

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
