#/bin/bash

go build . 
mv app.yml app_dev.yml
./cisasntyi -c app_dev.yml
