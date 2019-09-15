#!/bin/bash
# docker="$DOCKER"/docker.exe
echo ""
which docker
docker build -t go-subway -f docker/Dockerfile .