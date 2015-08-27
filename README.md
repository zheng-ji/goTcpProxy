[![Build Status](https://travis-ci.org/zheng-ji/goTcpProxy.svg)](https://travis-ci.org/zheng-ji/goTcpProxy)

## goTcpProxy 

A Tcp Proxy Server Written By Go

### Description

* support tcp proxy
* support catching exit signal 
* support customize Your connect Params such as `max_avail_conn`, `ttl`

### How To Compile

```
cd $GOPATH;
git clone http://github.com/zheng-ji/goTcpProxy;
cd src;
make
```

### How To Use

```
Usage of ./goTcpProxy:
    -dest="172.16.1.250:3306": where proxy server forward requests to.
    -max_avail_conn=25: max active connection.
    -max_wait_conn=10000: max connections in the queue wait for servers.
    -src="127.0.0.1:8999": server's host.
    -ttl=20: timeout of read and write
```

```
./goTcpProxy -src="127.0.0.1:8999" -dest="172.16.1.250:3306"
```

----
MIT LICENSE

