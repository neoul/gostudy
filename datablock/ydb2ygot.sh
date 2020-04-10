#!/bin/bash

cd model
go generate

# ydb -r pub -a uss://test -d -f demo.yaml &
# PUBPID=$!

# go run demo.go
# kill $PUBPID
# cd -
