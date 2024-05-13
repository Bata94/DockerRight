#! /usr/bin/env bash

docker build --target prod -t ghcr.io/bata94/dockerright:latest .
docker push ghcr.io/bata94/dockerright:latest
