# Go Backend Development - goserve

## API Framework Design
![Request-Response-Design](.extra/docs/api-structure.png)

## API DOC
[![API Documentation](https://img.shields.io/badge/API%20Documentation-View%20Here-blue?style=for-the-badge)](https://documenter.getpostman.com/view/1552895/2sA3XWdefu)

# Installation Instruction
vscode is the recommended editor - dark theme 

### Get the repo 

```bash
git clone https://github.com/unusualcodeorg/goserve.git
```

### Run Docker Compose
- Install Docker and Docker Compose. [Find Instructions Here](https://docs.docker.com/install/).

```bash
docker-compose up --build
```
-  You will be able to access the api from http://localhost:8080

### Run Tests
```bash
TODO
```

If having any issue
- Make sure 8080 port is not occupied else change SERVER_PORT in **.env** file.
- Make sure 27017 port is not occupied else change DB_PORT in **.env** file.
- Make sure 6379 port is not occupied else change REDIS_PORT in **.env** file.

## Run on the local machine
Change the following hosts in the **.env**
- DB_HOST=localhost
- REDIS_HOST=localhost

Best way to run this project is to use the vscode `Run and Debug` button. Scripts are available for debugging and template generation on vscode.

```bash
go mod tidy
```
or

```bash
make install
```

### Running the app
```bash
go run cmd/main.go
```

or
```bash
make run
```


## Find this project useful ? :heart:
* Support it by clicking the :star: button on the upper right of this page. :v:

## More on YouTube channel - Unusual Code
Subscribe to the YouTube channel `UnusualCode` for understanding the concepts used in this project:

[![YouTube](https://img.shields.io/badge/YouTube-Subscribe-red?style=for-the-badge&logo=youtube&logoColor=white)](https://www.youtube.com/@unusualcode)

## Contribution
Please feel free to fork it and open a PR.