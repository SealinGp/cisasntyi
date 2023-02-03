#!/bin/bash

go mod tidy
GOOS=linux GOARCH=amd64 go build -o cisasntyi-binary .

if [ -d "cisasntyi" ]; then
  rm -rf cisasntyi
fi

mkdir cisasntyi
mv cisasntyi-binary  cisasntyi

if [ ! -f "app_dev.yml" ] ; then
  cp app.yml app_dev.yml
fi

cp app_dev.yml cisasntyi
cp start.sh cisasntyi

tar -cvzf cisasntyi-linux-amd64.tar.gz cisasntyi && rm -rf cisasntyi
