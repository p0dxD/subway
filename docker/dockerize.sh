#!/bin/bash
# docker="$DOCKER"/docker.exe
echo "Creating image"
docker build -t go-subway -f docker/Dockerfile .
echo "Cealing old images"
yes | docker image prune

echo "Stopping old container"
old_container=$(docker ps -a | grep "go-subway" | awk '{ print $1 }')
if [ ! -z $old_container ]; then
    docker stop $old_container
    docker rm  $old_container
fi

echo "Starting new container"
docker run -d -it -p 3001:3001 go-subway