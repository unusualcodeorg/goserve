# goserve - Go Backend Architecture

[![Docker Compose CI](https://github.com/unusualcodeorg/goserve/actions/workflows/docker_compose.yml/badge.svg)](https://github.com/unusualcodeorg/goserve/actions/workflows/docker_compose.yml)

![Banner](.extra/docs/goserve-banner.png)

## Create A Blog Service 

This project is a fully production-ready solution designed to implement best practices for building performant and secure backend REST API services. It provides a robust architectural framework to ensure consistency and maintain high code quality. The architecture emphasizes feature separation, facilitating easier unit and integration testing.

# Framework
- Go
- Gin
- jwt
- mongodriver
- go-redis
- Validator
- Viper
- Crypto

## Highlights
- API key support
- Token based Authentication
- Role based Authorization
- Unit Tests
- Integration Tests
- Modular codebase

# Architecture
The goal is to make each API independent from one another and only share services among them. This will make code reusable and reduce conflicts while working in a team. 

The APIs will have separate directory based on the endpoint. Example `blog` and `blogs` will have seperate directory whereas `blog`, `blog/author`, and `blog/editor` will share common resources and will live inside same directory.

## Startup Flow
cmd/main → startup/server → module, mongo, redis, router → api/[feature]/middlewares → api/[feature]/controller -> api/[feature]/service, authentication, authorization → handlers → sender

## API Structure
```
Sample API
├── dto
│   └── create_sample.go
├── model
│   └── sample.go
├── controller.go
└── service.go
```

- Each feature API lives under `api` directory
- The request and response body is sent in the form of a DTO (Data Transfer Object) inside `dto` directory
- The database collection model lives inside `model` directory
- Controller is responsible for defining endpoints and corresponding handlers
- Service is the main logic component and handles data. Controller interact with a service to process a request. A service can also interact with other services.
 
# Project Directories
1. **api**: APIs code 
2. **arch**: It provide framework and base implementation for creating the architecture
3. **cmd**: main function to start the program
4. **common**: code to be used in all the apis
5. **config**: load environment variables
6. **keys**: stores server pem files for token
7. **startup**: creates server and initializes database, redis, and router
8. **tests**: holds the integration tests
9. **utils**: contains utility functions

## Helper/Optional Directories
1. **.extra**: mongo script for initialization inside docker, other web assets and documents
2. **.github**: CI for tests
3. **.tools**: api code, RSA key generator, and .env copier
4. **.vscode**: editor config and debug launch settings

# Documentation
Check the [Wiki](https://github.com/unusualcodeorg/goserve/wiki/Architecture) for the detailed documentation on the architecture. 

Note: We will keep on adding documentations progressively

## API Design
![Request-Response-Design](.extra/docs/api-structure.png)

## API DOC
[![API Documentation](https://img.shields.io/badge/API%20Documentation-View%20Here-blue?style=for-the-badge)](https://documenter.getpostman.com/view/1552895/2sA3XWdefu)

# Installation Instruction
vscode is the recommended editor - dark theme 

### 1. Get the repo 

```bash
git clone https://github.com/unusualcodeorg/goserve.git
```

### 2. Generate RSA Keys
```
go run .tools/rsa/keygen.go
```

### 3. Create .env files
```
go run .tools/copy/envs.go 
```

### 4. Run Docker Compose
- Install Docker and Docker Compose. [Find Instructions Here](https://docs.docker.com/install/).

```bash
docker-compose up --build
```
-  You will be able to access the api from http://localhost:8080

### 5. Run Tests
```bash
docker exec -t goserver go test -v ./...
```

If having any issue
- Make sure 8080 port is not occupied else change SERVER_PORT in **.env** file.
- Make sure 27017 port is not occupied else change DB_PORT in **.env** file.
- Make sure 6379 port is not occupied else change REDIS_PORT in **.env** file.

# Run on the local machine
```bash
go mod tidy
```

Keep the docker container for `mongo` and `redis` running and **stop** the `goserve` docker container

Change the following hosts in the **.env** and **.test.env**
- DB_HOST=localhost
- REDIS_HOST=localhost

Best way to run this project is to use the vscode `Run and Debug` button. Scripts are available for debugging and template generation on vscode.

### Optional - Running the app from terminal
```bash
go run cmd/main.go
```

# Template
New api creation can be done using command. `go run .tools/apigen.go [feature_name]`. This will create all the required skeleton files inside the directory api/[feature_name]

```bash
go run .tools/apigen.go sample
```

## Find this project useful ? :heart:
* Support it by clicking the :star: button on the upper right of this page. :v:

## More on YouTube channel - Unusual Code
Subscribe to the YouTube channel `UnusualCode` for understanding the concepts used in this project:

[![YouTube](https://img.shields.io/badge/YouTube-Subscribe-red?style=for-the-badge&logo=youtube&logoColor=white)](https://www.youtube.com/@unusualcode)

## Contribution
Please feel free to fork it and open a PR.