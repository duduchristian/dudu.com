#!/bin/zsh
set -e

docker stop my_server || true
sleep 20
./deploy.sh -use-tuner
docker stop my_server || true
sleep 20
./deploy.sh