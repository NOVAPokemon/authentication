#!/bin/bash

set -e

env GOOS=linux GOARCH=amd64 go build -a -v -o executable .
docker build -t brunoanjos/authentication-test:latest .
docker push brunoanjos/authentication-test:latest
