# financial-api ![develop](https://github.com/lotostudio/financial-api/actions/workflows/develop.yml/badge.svg)

Backend service for financial application

## Variables

Use these variables to run project in `.env` file

```dotenv
LOG_LEVEL=INFO

GIN_MODE=release    # For prod

DB_HOST=<host>
DB_PORT=<port
DB_NAME=<database>
DB_USER=<username>
DB_PASSWORD=<password>
DB_SSL_MODE=<mode>

AUTH_ACCESS_TOKEN_TTL=<ttl>
AUTH_REFRESH_TOKEN_TTL=<ttl>
AUTH_REFRESH_TOKEN_LENGTH=<length>
AUTH_PASSWORD_SALT=<salt>
AUTH_JWT_KEY=<key>

ACCOUNT_CARD_CASH_LIMIT=<limit>
ACCOUNT_LOAN_DEPOSIT_LIMIT=<limit>
```

## Commands

`go generate` - generate mock classes _(in package)_

`make fmt` - format whole project with gofmt _(do it before any commit)_

`make swag` - generate openapi documentation

`make cover` - run unit tests and show coverage report

`make build` - build project

`make run` - build and run project

## Docker

Use dockerfiles in `build` directory for building images and running containers

Use `build/Dockerfile` for building images on unix systems.
Use `build/Dockerfile.multi` for building images on non-unix systems
