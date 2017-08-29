#!/bin/bash

# To build for the Docker image
# env GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -v server.go

# To build for OSX
go build -o server server.go

