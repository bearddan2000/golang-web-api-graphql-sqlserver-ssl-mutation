#!/usr/bin/env bash

# Set $GOPATH so the dependcies will be found
# dependcies are expected in a src dir
export GOPATH="`pwd`/src"

mkdir src

cd src

# Command `go get` calls the underlying command `git clone`
apt-get -y install git

# List of dependcies
go get github.com/graphql-go/graphql \
&& go get github.com/denisenkom/go-mssqldb

cd ../

go run ./main.go > /tmp/mylog/out
