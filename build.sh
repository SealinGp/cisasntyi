#!/bin/bash

go mod tidy
go build -x .

tar -cvzf apple.tar.gz app.yml apple start.sh