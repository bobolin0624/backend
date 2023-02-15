# Backend

## Setup Development Environment

1. [Install Docker Desktop](https://www.docker.com/get-started/)
1. [Download Go 1.19](https://go.dev/dl/)
1. Clone the repo

```sh
$ git clone git@github.com:taiwan-voting-guide/backend.git
$ cd backend
```

1. Setup postgres

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

1. Start the server

```sh
$ go run main.go
```

1. Copy env file

```sh
$ cp .env.example .env
```

1. Init politicians testing

```sh
./scripts/init_politicians.sh
```

1. (optional) [Download pgAdmin](https://www.pgadmin.org/download/)
