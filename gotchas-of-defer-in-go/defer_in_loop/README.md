
Step 1. Start a postgres server locally and set the server to allow at most 4 connections. 
```sh
docker run --rm --name=pg1 -p 15432:5432 -e POSTGRES_PASSWORD=password postgres -c max_connections=4
```

Step 2. 
```sh
go run main.go
```
