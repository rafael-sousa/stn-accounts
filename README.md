[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/rafael-sousa/stn-accounts)](https://github.com/rafael-sousa/stn-accounts)
[![Go Report Card](https://goreportcard.com/badge/github.com/rafael-sousa/stn-accounts)](https://goreportcard.com/report/github.com/rafael-sousa/stn-accounts)
[![Go Reference](https://pkg.go.dev/badge/github.com/rafael-sousa/stn-accounts.svg)](https://pkg.go.dev/github.com/rafael-sousa/stn-accounts)

<p align="center">
  <h3 align="center">STN Accounts</h3>

  <p align="center">
    A small REST API written in Go(lang)
  </p>
</p>

<summary><h2 style="display: inline-block">Table of Contents</h2></summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
        <li><a href="#api-routes">API Routes</a></li>
      </ul>
    </li>
    <li>
        <a href="#development">Development</a>
        <ul>
            <li><a href="#application-layers">Application Layers</a></li>
            <li><a href="#project-structure">Project Structure</a></li>
            <li><a href="#testing">Testing</a></li>
            <li><a href="#environment-config">Environment Config</a></li>
            <li><a href="#dependencies">Dependencies</a></li>
            <li><a href="#troubleshooting">Troubleshooting</a></li>
        </ul>
    </li>
    <li><a href="#license">License</a></li>
    <li><a href="#acknowledgements">Acknowledgements</a></li>
  </ol>


## About The Project

This project aims to fulfill the proposed technical challenge applying software development patterns and following the Go's best practices and conventions. The API exposes endpoints that handles operations on `Account` and `Transfer` core domain of a digital bank.

## Getting Started

This section describes the steps to get a local copy up and running.

### Prerequisites

* [Go 1.15+](https://golang.org/)
* [Docker 19+](https://www.docker.com/)
* [Docker Compose 3.3+](https://docs.docker.com/compose/)

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/rafael-sousa/stn-accounts.git
   ```

2. Run the API within a docker container
   ```sh
   make start
   ```


3. Once running, the swagger UI is available at: [localhost:3000/swagger/index.html](http://localhost:3000/swagger/index.html)
   

4. Stop the running API
   ```sh
   make stop
   ```


### API Routes

The following table shows the current available endpoints

| METHOD | PATH                           | AUTH |
|--------|--------------------------------|------|
| GET    | /accounts                      |      |
| GET    | /accounts/{id}/balance         |      |
| POST   | /accounts                      |      |
| POST   | /login                         |      |
| GET    | /transfers                     | X    |
| POST   | /transfers                     | X    |

## Development

This section portrays the application architecture and how their elements are laid

### Application Layers

The following image displays the different conceptual layers:

![](https://user-images.githubusercontent.com/12838206/110670269-fa14f400-81ab-11eb-9006-aa9aeda1434d.PNG)

### Project Structure

The next graph shows the application folder layout along with its short description:

```
├───cmd
│   └───stn-accounts         ; holds the application entry point
├───docs                     ; keeps OpenAPI resource files
└───pkg                      ; api source code
    ├───controller
    │   └───rest             ; maintains code related to REST endpoints
    │       ├───body         ; request and response body models
    │       ├───jwt          ; handle jwt creation and parsing 
    │       ├───middleware   ; custom api middlewares
    │       ├───response     ; standard response writer functions
    │       └───routing      ; routes exposed by the api
    ├───model
    │   ├───dto              ; transfer data structs between different layers
    │   ├───entity           ; database models
    │   ├───env              ; environment models
    │   └───types            ; custom application types
    ├───repository
    │   └───mysql            ; mysql repository implementation
    │       └───migrations   ; mysql-specific migration files
    ├───service
    │   └───validation       ; maintains complex business rules for reuse
    └───testutil             ; centralize test utilities
```



### Testing

In order the run the following commands, a go installation is required with a version 1.15+

1. Running the application tests
   ```sh
   make test
   ```

2. Format and analyze source code
   ```sh
   make lint
   ```

### Environment Config

The application can be configured overrinding the following environment variables:

| NAME                 | TYPE   | DESCRIPTION                                  | DEFAULT VALUE    |
|----------------------|--------|----------------------------------------------|------------------|
| DB_PORT              | UINT   | Database connection port                     | 3306             |
| DB_USER              | STRING | Database user name                           | admin            |
| DB_PW                | STRING | Database user password                       | admin            |
| DB_HOST              | STRING | Database user password                       | localhost        |
| DB_NAME              | STRING | Database name                                | stn_accounts     |
| DB_DRIVER            | STRING | Database driver                              | mysql            |
| DB_MAX_OPEN_CONNS    | UINT   | Maximum open connection number               | 10               |
| DB_MAX_IDLE_CONNS    | UINT   | Maximum idle connection number               | 10               |
| DB_CONN_MAX_LIFETIME | UINT   | Maximum connection lifetime                  | 0                |
| DB_PARSE_TIME        | BOOL   | Database flag for parsing time automatically | true             |
| PORT                 | UINT   | Http server port                             | 3000             |
| JWT_SECRET           | STRING | Secret used to generate and parse JWT Tokens | rest-app@@secret |
| JWT_EXP_TIMEOUT      | UINT   | JWT Token timeout in minutes                 | 30               |

### Dependencies
The following table lists the direct dependencies used by the application. A complete list can be found on [go.mod file](https://github.com/rafael-sousa/stn-accounts/blob/main/go.mod)

| NAME                                              | VERSION | DESCRIPTION                                |
|---------------------------------------------------|---------|--------------------------------------------|
| [jwt-go](github.com/dgrijalva/jwt-go)             | v3.2.0  | Used for generating and parsing jwt tokens |
| [chi](github.com/go-chi/chi)                      | v4.0.2  | Provides routes and http middlewares       |
| [mysql](github.com/go-sql-driver/mysql)           | v1.5.0  | Database driver                            |
| [migrate](github.com/golang-migrate/migrate)      | v3.5.4  | Migration tool                             |
| [dockertest](github.com/ory/dockertest/v3)        | v3.6.3  | Testing tool for running repository tests  |
| [zerolog](github.com/rs/zerolog)                  | v1.20.0 | Application logger                         |
| [go-envconfig](github.com/sethvargo/go-envconfig) | v0.3.2  | Environment config parser                  |
| [http-swagger](github.com/swaggo/http-swagger)    | v1.0.0  | OpenAPI implementation                     |
| [swag](github.com/swaggo/swag)                    | v1.7.0  | Static swagger files generator             |
| [crypto](golang.org/x/crypto)                     | v0.0.0  | Password encrypter                         |

## Troubleshooting

* Error when mounting docker-compose data volume
    1. Edit the docker agent settings
    2. On the left menu, navigate to Resources -> File Sharing
    3. Click at the '+' button and add the download repository directory
    4. Apply the settings and wait the service restart

* Error during database test execution "Could not start resource ... No connection could be made because the target machine actively refused it"
    1. Add environment variable DOCKER_HOST a with value `tcp://127.0.0.1:2375`
    2. Edit the docker agent settings
    3. On the left menu, navigate to General
    4. Check the "Expose daemon on tcp://localhost:2375 without TLS" option

* Cleaning up docker containers and volume
    1. Stop docker-compose executions: docker-compose down
    2. Delete all containers using the following command: docker rm -f $(docker ps -a -q)
    3. Delete all volumes using the following command:docker volume rm $(docker volume ls -q)


## License

Distributed under the MIT License. See `LICENSE` for more information.

## Acknowledgements

1. [Complete Guide to Create Docker Container for Your Golang Application](https://levelup.gitconnected.com/complete-guide-to-create-docker-container-for-your-golang-application-80f3fb59a15e)
2. [Slimming Down Your Docker Images](https://towardsdatascience.com/slimming-down-your-docker-images-275f0ca9337e)
3. [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
4. [Unit Testing made easy in Go](https://medium.com/rungo/unit-testing-made-easy-in-go-25077669318)
5. [Golang's Mocking Techniques](https://www.youtube.com/watch?v=LEnXBueFBzk)
6. [Go: Are pointers a performance optimization?](https://medium.com/@vCabbage/go-are-pointers-a-performance-optimization-a95840d3ef85)
7. [When to use pointers in Go](https://medium.com/@meeusdylan/when-to-use-pointers-in-go-44c15fe04eac)