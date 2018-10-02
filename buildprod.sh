#! /bin/bash

# Build web and other services

cd ~/Go_WS/src/video_server/api
env GOOS=linux GOARCH=amd64 go build -o ../bin/api

cd ~/Go_WS/src/video_server/scheduler
env GOOS=linux GOARCH=amd64 go build -o ../bin/scheduler

cd ~/Go_WS/src/video_server/streamserver
env GOOS=linux GOARCH=amd64 go build -o ../bin/streamserver

cd ~/Go_WS/src/video_server/web
env GOOS=linux GOARCH=amd64 go build -o ../bin/web

echo 全部编译完成