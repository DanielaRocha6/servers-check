# servers-check

## Requirements

This project requires Go, yarn and CockroachDB. They can be installed in
* [Go](https://golang.org/dl/)
* [yarn](https://classic.yarnpkg.com/en/docs/install)
* [CockroachDB](https://www.cockroachlabs.com/docs/stable/install-cockroachdb.html)


## DB local setup

#### Start CockroachDB
```
cockroach start-single-node --insecure
```

#### Create database
```
cockroach sql --insecure
create database domains;
```
> *Note:* It will deploy the database in `http://localhost:26257`and a Dashboard of it in `http://localhost:8080`

## Project local setup

### Back setup
Go to `back` directory

For installing dependencies and building project, run:
```
make
```

Run project, which will run in `http://localhost:8090`
```
./back
```

### Front setup

Go to `front` directory

#### Install dependencies
```
yarn install
```

#### Compiles and hot-reloads for development
```
yarn serve
```
