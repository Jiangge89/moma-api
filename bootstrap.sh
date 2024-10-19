#!bin/bash 

export GODEBUG=netdns=cgo

exec ./bin/api 1>>./log/api.error 2>&1
