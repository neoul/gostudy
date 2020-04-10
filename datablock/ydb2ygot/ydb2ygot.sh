#!/bin/bash

# Change ygot branch to avoid the code error.
go get -u github.com/openconfig/ygot
git checkout 6daf745bd5f14eda714e98cec83884e5b3954898

cd model
go generate

# ydb -r pub -a uss://test -d -f demo.yaml &
# PUBPID=$!

# go run demo.go
# kill $PUBPID
# cd -
