#!/bin/bash

go mod tidy
go build -x .

tar -cvzf cisasntyi.tar.gz app.yml cisasntyi start.sh