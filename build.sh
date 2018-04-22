#!/bin/sh

docker build -t=donutloop/stories:v1.0.0 -f=Dockerfile.stories .
docker build -t=donutloop/frontend:v1.0.0 -f=Dockerfile.front .
