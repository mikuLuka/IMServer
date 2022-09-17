# 一个简单的即时通讯软件

## 1.服务端

server.go ,创建基础server

user.go, server端对用户的操作

广播用户上线下线，查询用户，更新用户名，超时强踢



## 2.客户端

client.go， 客户端连接

clientMenuFunc.go，客户端功能

公聊模式，私聊模式，更新用户名



## 3.使用

Windows环境下，

服务端：使用 `go build server.go user.go main.go`  创建server.exe

客户端：使用  `go build client.go clientMenuFunc.go`  创建client.exe