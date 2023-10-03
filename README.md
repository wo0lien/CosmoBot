# CosmoBot

## Generate SDK

`oapi-codegen -generate "types,client" -package api -o ./internal/api/client.gen.go ./swagger.json`

## Using the program

This program uses go-sqlite3 package that requires CGO and is therefore recommnended to run under Linux-UNIX environment.
