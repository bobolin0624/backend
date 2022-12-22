# Backend

## Setup Development Environment
1. [Install Docker Desktop](https://www.docker.com/get-started/)
2. [Download Go 1.19](https://go.dev/dl/)
3.
```sh
$ git clone git@github.com:taiwan-voting-guide/backend.git

# setup postgres
$ docker run --name postgres -e POSTGRES_PASSWORD=password -d postgres

# start the server
$ go run main.go 
```
4. (optional) [Download pgAdmin](https://www.pgadmin.org/download/)
