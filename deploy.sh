#!/bin/zsh
set -e

docker stop my_server || true
docker rm my_server || true
docker rmi testing:1.0 || true
GOOS=linux go build cmd/main.go
docker build . -t testing:1.0
docker run -p 8080:8080 --name=my_server -d testing:1.0 -use-fasthttp