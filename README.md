# SERVICES API

You can find DB schema here - https://dbdiagram.io/d/6456768fdca9fb07c49eceb3

I decided to use PostgreSQL as main database and Gin framework for simplicity and
functionality that we can use out of the box.

Things to improve: 
- cover by unit tests 4xx and 5xx errors
- user friendly validation errors

## Available endpoints
`GET /v1/services` - returns services list with filtering, sorting and pagination

`GET /v1/services/:id/versions` - get all versions for particular service

### Project insights

Covered by unit tests 75%.

#### Main libraries and reasons why

`gin` - http framework

`sqlc` - library to generate go code from sql queries

`mockgen` - to generate mocks for unit tests

`migrate` - to run migrations on database

`zerolog` - for logging (basic things)

`viper` - for configuration

`testify` - for unit tests

## How to run locally

In order to run it locally you should have `migrate`, `sqlc`, `mockgen` locally

1. Run postgres in docker - `make postgres`
2. Create DB - `make createdb`
3. Run migrations - `make migrateup`
4. Run http server - `make server`

P.s. check `Makefile` for more commands