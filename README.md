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
    -c=10000: max connections in the queue wait for servers.
    -dest="172.16.1.250:3306": where proxy server forward requests to.
    -n=25: max active connection.
    -src="127.0.0.1:8999": proxy server's host.
    -ttl=20: timeout of read and write
```

```
./goTcpProxy -src="127.0.0.1:8999" -dest="172.16.1.250:3306"
```

License
-------

Copyright (c) 2015 released under a MIT style license.

