# Servers check!
This project lets get information about domains and its servers.
This is the retrieved information:
* Title
* Logo
* Status - `up or down`
* SSL grade - `lowest SSL grade among its servers`
* Previous SSL grade - `SSL grade the domain had past one hour or more`
* Servers changed - `If servers had changed past one hour or more`
* Servers - Retrieved from [SSL Labs](https://www.ssllabs.com/projects/ssllabs-apis/index.html)
    * Address
    * SSL grade
    * Country
    * Owner

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
> In `http://localhost:8090/checkDomain/:domain` it can be consulted a domain 
>
> In `http://localhost:8090/allDomains`it can be consulted all the past searched domains.


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
> It will be deployed in `http://localhost:8081`
