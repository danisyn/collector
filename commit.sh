#!/bin/bash

git add .

read -p "Commit message: " message

git commit -m "$message"

git push origin main

cp dockerfiles/Dockerfile.latest Dockerfile

docker build --tag daniels7/collector:latest .

docker push daniels7/collector:latest

cp dockerfiles/Dockerfile.1.20-alpine Dockerfile

docker build --tag daniels7/collector:1.20-alpine .

docker push daniels7/collector:1.20-alpine
