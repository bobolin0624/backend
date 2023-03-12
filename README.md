# Backend

## Setup Development Environment

1. [Install Docker Desktop](https://www.docker.com/get-started/)
1. [Download Go 1.19](https://go.dev/dl/)
1. Clone the repo

```sh
git clone git@github.com:taiwan-voting-guide/backend.git
cd backend
```

4. Setup postgres

```sh
docker run \
-v `pwd`/init.sql:/docker-entrypoint-initdb.d/init.sql:ro \
--name pg \
-e POSTGRES_USER=backend_user \
-e POSTGRES_PASSWORD=password \
-e POSTGRES_DB=tvg \
-p 5432:5432 \
-d postgres
```

5. Copy env file

```sh
cp .env.example .env
# copy google client id from slack #tech
```

6. Start the server

```sh
go run main.go
```

7. Init politicians testing

```sh
./scripts/init_politicians.sh
```

1. (optional) [Download pgAdmin](https://www.pgadmin.org/download/)

## Staging APIs

### Create a staging record
`table`: This refers to the target table to be processed.

`searchBy`: These are the fields used to search the target table. If the data already exists in the target table, a create record is created. If not, an update record is created.

`fields`: These are the individual fields used for the target table. If a field is a reference to another table, the reference's ID is searched for and replaced.

```
POST /workspace/staging

{
    "table": "parties",
    "searchBy": {
        "name": "民主進步黨" 
    },
    "fields": {
        "name": "民主進步黨",
        "chairman": "賴清德",
        "established_date": "2012-12-12",
        "filing_date date": "2012-12-12",
        "main_office_address": "kkkkkkkk",
        "mailing_address": "aaaaaaaaaa",
        "phone_number": "091123321",
        "status": "三小"
    }
}
```
```
POST /workspace/staging

{
    "table": "politicians",
    "searchBy": {
        "name": "許淑華",
        "birthdate": "1975-05-22"
    },
    "fields": {
        "name": "許淑華",
        "birthdate": "1975-05-22",
        "sex": "female",
        "current_party_id": {
            "table": "parties",
            "searchBy": {
                "name": "中國國民黨"
            }
        }
    }
}
```

## Troubleshoot

1. _Wired db connection/schema error_: Try pulling the latest master and re-init pg.
