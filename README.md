# Backend

## Setup Development Environment
1. [Install Docker Desktop](https://www.docker.com/get-started/)
2. [Download Go 1.19](https://go.dev/dl/)
3. clone the repo
```sh
$ git clone git@github.com:taiwan-voting-guide/backend.git
```
4. setup postgres
```sh
$ docker run --name postgres -e POSTGRES_PASSWORD=password -d postgres
```
5. start the server
```sh
$ go run main.go 
```
6. (optional) [Download pgAdmin](https://www.pgadmin.org/download/)
