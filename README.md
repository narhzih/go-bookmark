# Requirements 
To be able to run this bookmarking Api locally, you'll need the following installed on your computer: 
### Go (Version 1.18)
This installation is quite optional, but you'll need it if you don't want to start up the server with docker
You can find instructions on how to install Go on your machine [here](https://go.dev/doc/install)
### Docker
If you don't already have docker installed, you can also find instructions on how to install docker on your machine [here](https://docs.docker.com/desktop/)
### Postgres
Use [this](https://www.postgresqltutorial.com/postgresql-getting-started/install-postgresql-macos/) link for instructions how to install postgres on your machine, take note of the user you created, the password and the name of 
the database, you'll use them as values for the environment variables

### Make 
To install `make` on MacOS, you can use homebrew as follows: 
```bash
$ brew install make
```
You can check online resources for other ways on how to install `make` on Linux and Windows machine

### Golang Migrate
To install golang-migrate on MacOS, you can use homebrew as follows:
```bash
$ brew install golang-migrate
```
You can check online resources for other ways on how to install `golang-migrate` on Linux and Windows machine

After installing all required software, run the following command to create a copy of .env.example

```bash
$ cp .env.example .env
```

### Environment variables

- `APP_PORT`: The port on which the API server should run.

- `POSTGRES_DB_HOST`: Name/IP of the PostgreSQL server. This defaults to postgres on dev environment
- `POSTGRES_DB_PORT`: Port number that postgres runs on your local environment (defaults to 5432)

Other DB variables supplied like `POSTGRES_USER, POSTGRES_PASSWORD and POSTGRES_DB` are particular to your local environment.

### Start up the server
Run `docker-compose up --build` in the project root directory to start the docker processes

## Migrations
After installing `make` and `golang migrate`, run the following command to create the database tables for your database:
```bash
$ make migrate-up
```

## Optionally setting up without docker 
If you're setting up the local environment without docker (although, docker is highly recommended), you'll need the following 
command to start the server on your local machine in this project root folder like so: 
```bash
$ go run main.go
```
Note: I assume that you already have Go (version 1.18) already installed on your computer
