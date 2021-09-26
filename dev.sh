#/bin/bash
echo "开始构建..."
go build .
if [ $? -ne 0 ] ; then
  echo "构建失败,请检查您的go环境是否正常"
  exit
fi
echo "构建成功,准备运行..."


if [ ! -f "app_dev.yml" ] ; then
  echo "检查到您是首次运行,开始生成配置文件..."
  cp app.yml app_dev.yml  

  if [ $? -ne 0 ] ; then 
    echo "生成配置文件失败."
    exit
  fi

  echo "生成配置app_dev.yml文件成功,请修改该文件配置后重新运行"
  exit
fi

./cisasntyi -c app_dev.yml
