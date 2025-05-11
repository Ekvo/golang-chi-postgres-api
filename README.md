# golang-chi-postgresql ~> REST API :)

This code demonstrates the CRUD principle in implementing a REST API using golang/chi and PostgresSQL.  

[Documentation link - ./docs/docs.go](https://github.com/Ekvo/golang-chi-postgres-api/tree/main/docs/docs.go "https://github.com/Ekvo/golang-chi-postgres-api/tree/main/docs/docs.go")
### Directory structure
```
.
├── cmd/app
│   └──── main.go     
├── docs  
│   └──── docs.go         // documentation
├── internal
|   ├── config
|   │   └──── config.go   
|   ├── model
|   │   └──── model.go    // data models define
|   ├── server  
|   │   └──── server.go   // init for http.Server
|   ├── servises           
|   │   ├── serializer.go // response computing & format
|   │   └── validator.go  // json checker        
|   ├── source
|   │   ├── query.go      // SQL query for model
|   │   └── source.go     // init for *sql.DB
|   ├── transport 
|   │   ├── middlweare.go    
|   │   ├── route.go      
|   │   └── transport.go  // router binding
|   └── variables.go      
|       └──── variables.go  // only const, var
├── pkg/common 
│   └──── common.go         // tools function
├── ...
 .env
 compose.yaml
 Dockerfile
...
```
#### * [Golang - instal](https://go.dev/doc/install "https://go.dev/doc/install")
This code written in golang v1.21.0

#### * A pure Go postgres driver for Go's database/sql  
```bash
go get github.com/lib/pq
```
#### Small and composable router -> [chi](https://pkg.go.dev/github.com/go-chi/chi "https://pkg.go.dev/github.com/go-chi/chi")
```bash
go get github.com/go-chi/chi/v5
```
#### Config use the package(s):  
Local run - Load **.env** file to initialize database and http.Server propirties, in Docker Image use ENV varialbe see Dockerfile
```bash
go get github.com/joho/godotenv
```

[Viper](https://github.com/spf13/viper "https://github.com/spf13/viper") for read .env file and ENV variavles
```bash
go get github.com/spf13/viper
```

#### * Docker
Have *compose.yaml* -> from root project
1. postgres:alpine
2. golang:1.21.0

see *Dockerfile* and *compose.yaml*
```bash
docker compose up
```

#### * Exist .golangci.yaml
```bash
# install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```
```bash
# start from main folder
golangci-lint run
```

#### * Test 

Simple to understand and use - testing tools
```bash
go get github.com/stretchr/testify
```

From root project
```bash
go test ./...
```

for more details _**use** From package_test.go_ 
```bash
go test . -coverprofile=coverage.out
```
then use for see more details
```bash
go tool cover -html=coverage
```

##### Сoverage of packages

| path to package                   | percent *%* |
|:----------------------------------|:-----------:|
| ./internal/source/source.go       |    58.3     |
| ./internal/source/query.go        |    66.7     |
|                                   |             |
| ./internal/transport/midlweare.go |    88.9     |
| ./internal/transport/route.go     |    68.5     |
| ./internal/transport/transport.go |    100.0    |
|                                   |             |
| ./internal/common/common,go       |    90.5     |

#### * CURL
 1. Create new task

```http request
curl -X POST -H "Content-Type: application/json" -d '{"task_update":{"description":"test", "note":"test"}}' http://127.0.0.1:3000/task/
```
 2. Read created task

```http request
curl -i -H "Accept: application/json" http://127.0.0.1:3000/task/desc/1/0
```

*Thank you for your time:)*  
ps. to use *locally* you need to change `HOST` in `.env` to `127.0.0.1`



