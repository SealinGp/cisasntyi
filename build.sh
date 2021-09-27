#!/bin/bash

go mod tidy
go build -x -o cisasntyi-binary .

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

tar -cvzf cisasntyi.tar.gz cisasntyi && rm -rf cisasntyi
