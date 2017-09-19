#!/bin/bash

rootloc="`pwd`"

echo "building binary of rlog.."
cd $rootloc
go build -o rlog

echo "building rlog image"
docker build -t djavorszky/rlog .

echo "stopping previous version"
docker stop rlog-server
docker rm rlog-server

echo "starting rlog"
docker run -dit -p 1338:1338 --net=host --name rlog-server -v /home/javdaniel/go/src/github.com/djavorszky/ddn/dist/data:/log djavorszky/rlog:latest

echo "removing artefacts"
rm rlog
