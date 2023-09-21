#/bin/bash
echo "Start Building..."
go build .
if [ $? -ne 0 ] ; then
  echo "Buid failed, please check error message"
  exit
fi
echo "Build success,starting..."


if [ ! -f "app_dev.yml" ] ; then
  echo "Generate Config Files..."
  cp app.yml app_dev.yml  

  if [ $? -ne 0 ] ; then 
    echo "Faile to generate config file."
    exit
  fi

  echo "app_dev.yml is created, please midify this config for dev testing"
  exit
fi

./cisasntyi -c app_dev.yml
