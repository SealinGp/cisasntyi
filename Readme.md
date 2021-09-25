## for buy iphone 13 notification


cisasntyi: check iphone stock and send notification to your iphone

you need to install Bark app in you iphone 
change the app.yml config file 


### run local
```bash
bash dev.sh
```

### run in your private server
```bash
#build it
bash build.sh

#upload to your server
scp ./cisasntyi.tar.gz user@yourPrivateServerIp:/home

#run in your server
ssh -t user@ip cd /home
mkdir cisasntyi && mv cisasntyi.tar.gz cisasntyi && cd cisasntyi
bash start.sh
```