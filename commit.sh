#!/bin/bash

git add .

echo "Commit message: "

read message

git commmit -m $message

git push origin main

cp dockerfiles/Dockerfile.latest Dockerfile

docker build --tag daniels7/collector:latest .

docker push daniels7/collector:latest

cp dockerfiles/Dockerfile.1.20-alpine Dockerfile

docker build --tag daniels7/collector:1.20-alpine .

docker push daniels7/collector:1.20-alpine
