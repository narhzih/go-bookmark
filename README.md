# Requirements 
To be able to run the MyPipe API locally, you'll need the following installed on your computer: 
### Go (Version 1.18)
### Docker
### Postgres 
[Setting Up postgres on MacOS](https://www.postgresqltutorial.com/postgresql-getting-started/install-postgresql-macos/).
In following this link on how to install postgres, take note of the user you created, the password and the name of 
the database, you'll use them as values for the environment variables

### Make 
To install make on MacOS, you can use homebrew as follows: 
```bash
$ brew install make
```

### Golang Migrate
To install golang-migrate on MacOS, you can use homebrew as follows:
```bash
$ brew install golang-migrate
```


Run the below command to create a copy of .env.example

```bash
$ cp .env.example .env
```

### Environment variables

- `APP_PORT`: The port on which the API server should run.

- `POSTGRES_DB_HOST`: Name/IP of the PostgreSQL server. This defaults to postgres on dev environment
- `POSTGRES_DB_PORT`: Port number that postgres runs on your local environment (defaults to 5432)

Other DB variables supplied like `POSTGRES_USER, POSTGRES_PASSWORD and POSTGRES_DB` are particular to your local environment.

## Migrations
After installing make and golang migrate, run the following command to create the database tables for your database:
```bash
$ make migrate-up
```

## Authentication

## API documentation
