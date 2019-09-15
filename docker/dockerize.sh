#!/bin/bash
# docker="$DOCKER"/docker.exe
which docker
docker build -t go-subway -f docker/Dockerfile .