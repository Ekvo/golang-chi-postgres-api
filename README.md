# Golang/chi/PostgresQL REST API.

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
|   │   ├── route.go      // business logic
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
### [Golang - instal](https://go.dev/doc/install "https://go.dev/doc/install")

This code written in golang v1.21.0

Use package:  
 * Load **.env** file to initialize database and http.Server propirties
```
go get github.com/joho/godotenv
```
 *  A pure Go postgres driver for Go's database/sql  
```
go get github.com/lib/pq
```
 * Small and composable router -> [chi](https://pkg.go.dev/github.com/go-chi/chi "https://pkg.go.dev/github.com/go-chi/chi")
```
github.com/go-chi/chi/v5
```
 * Simple to understand and use - testing tools
```
github.com/stretchr/testify
```

### Docker
Have *compose.yaml* -> from root project
* golang:alpine
* postgres:alpine

see *Dockerfile* and *compose.yaml*

```bash
docker compose up
```

### Test 
From root project
```bash
go test ./...
```

./pkg/common         -> 94.1% of statements
./internal/source    -> 60.9% of statements
./internal/transport -> 70.1% of statements

_***Use** From package*_
```
go test . -coverprofile=coverage.out
go tool cover -html=coverage
```

### CURL
 * Create new task

```http request
curl -X POST -H "Content-Type: application/json" -d '{"task_update":{"description":"test", "note":"test"}}' http://127.0.0.1:3000/task/
```
 * Read created task

```http request
curl -i -H "Accept: application/json" http://127.0.0.1:3000/task/desc/1/0
```

*Thank you for you time:)*  
ps. to use *locally* you need to change `HOST` in `.env` to `127.0.0.1`



