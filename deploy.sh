#!/bin/zsh
set -e

count=$#
if [ $count -eq 0 ]; then
  echo "build fail"
  exit 1
fi


docker rmi testing:1.0 || true
GOOS=linux go build cmd/main.go
docker build . -t testing:1.0

server_num=$1
for i in $(seq "$server_num"); do
  server_name=my_server_$i
  docker stop "$server_name" || true
  docker rm "$server_name" || true
  docker run -m 8G --name="$server_name" -d testing:1.0 -use-fasthttp -use-tuner
done

#sleep 3
#go run cmd/test.go