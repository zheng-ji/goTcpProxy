[![Build Status](https://travis-ci.org/zheng-ji/goTcpProxy.svg)](https://travis-ci.org/zheng-ji/goTcpProxy)

## goTcpProxy 

A Tcp Proxy Server Written By Go

### Description

#### English
* A tcp proxy service
* Supprot multi backend severs 
* Consistent Hash Load Balance
* Auto detect down server, and remove it.
* Monitor backend health status

#### 中文

* TCP 代理服务
* 后端支持多个服务器
* 支持一致性哈希的负载均衡
* 自动检测失败的后端服务器，并移除
* 后端服务的健康检查接口

### How To Compile

```
cd $GOPATH;
git clone http://github.com/zheng-ji/goTcpProxy;
make
```

### How To Use

`config.yaml` 配置

```
bind: 0.0.0.0:9999      // 代理服务监听端口
wait_queue_len: 100     // 等待队列长度
max_conn: 10000         // 并发最大连接
timeout: 5              // 请求超时时间
failover: 3             // 后端服务允许失败次数 
stats: 0.0.0.0:19999    // 健康检查接口
backend:                // 后端服务列表
    - 127.0.0.1:80
    - 127.0.0.1:81
log:
    level: "info"
    path: "/Users/zj/proxy.log"
```

```
// 运行服务
./goTcpProxy -c=etc/config.yaml

// 查看后端服务健康状态
curl 127.0.0.1:19999/stats 
Server:127.0.0.1:80 FailTimes:9 isUp:false
Server:127.0.0.1:81 FailTimes:0 isUp:true
```

License
-------

Copyright (c) 2015 released under a MIT style license.
