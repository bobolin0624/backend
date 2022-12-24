# Backend

## Setup Development Environment
1. [Install Docker Desktop](https://www.docker.com/get-started/)
2. [Download Go 1.19](https://go.dev/dl/)
3. Clone the repo
```sh
$ git clone git@github.com:taiwan-voting-guide/backend.git
$ cd backend
```
4. Setup postgres
```sh
$ docker run \
-v `pwd`/init.sql:/docker-entrypoint-initdb.d/init.sql:ro \
--name pg \
-e POSTGRES_USER=backend_user \
-e POSTGRES_PASSWORD=password \
-e POSTGRES_DB=tvg \
-p 5432:5432 \
-d postgres
```
5. Start the server
```sh
$ go run main.go 
```
6. (optional) [Download pgAdmin](https://www.pgadmin.org/download/)
