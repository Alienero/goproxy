goproxy
=======
A simple proxy (two level proxy) writen in golang

## Install
You should have a golang env. Download[https://golang.org/dl/] 1.2 or later        
```
mkdir /home/goroxy
export GOAPTH=/home/goroxy
go get the github.com/Alienero/goproxy

go install github.com/Alienero/goproxy/client
go install github.com/Alienero/goproxy/server
go install github.com/Alienero/goproxy
cd $GOAPTH/bin
./key -host [your host]
# Instlled.
```
##Usage
```
cd $GOAPTH/bin

# Start the server
nohup ./server [OPinions] &
# Some flag
# the proxy server listen address.
#  -listen :8080
# the proxy server auth.
#  -password xxxx

# Start the client
./client [OPinions] 
# In the client ,you can firt start the program then config it
# Of cause you can use the args to config
#   -password xxx
# The Client listen address. Default is 127.0.0.1:808
#   -listen 127.0.0.1:808
# The remote server address(proxy server). Default is yim.so
#   -remote yim.so
# If you ues the windows,you want to hava the colorful print, you can set the iscolor true
# Default is false. It Only normal display in Cygwin Terminal
#   -iscolor
```
Hava Fun~Geaks
