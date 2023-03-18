# Notefan (Golang/Backend)

Notefan is Paypal Clone application built with Golang and Mysql

## Installation

install the dependencies

```sh
go get
go mod tidy
```

Migrate up Notefan required tables and seed data to tables

```sh
migrate -database "mysql://username:password@tcp(host:port)/database" -path database/migrations up
go run . seed
```

Migrate drop Notefan tables

```sh
migrate -database "mysql://username:password@tcp(host:port)/database" -path database/migrations drop
go run . seed
```

## License

MIT && NOTEFAN

**Open Source**
