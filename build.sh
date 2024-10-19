#!/bin/bash 

GIT_SHA=`git rev-parse --short HEAD || echo "NotGitVersion"` 
WHEN=`date '+%Y-%m-%d_%H:%M:%S'`

mkdir -p output/bin/ output/conf 

CGO_ENABLED=0 go build  -a -v -ldflags "-s -X main.GitSHA=${GIT_SHA} -X main.BuildTime=${WHEN}" -o output/bin/api 
chmod +x output/bin/api

cp bootstrap.sh output/
chmod +x output/bootstrap.sh 
