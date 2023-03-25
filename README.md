# go_redis

go_redis是一个用Go语言编写的简单Redis服务器，实现了一些基本的Redis命令和功能。它支持RESP协议，TCP连接，数据结构和数据库操作，以及AOF持久化。

## 安装

你可以使用go命令来安装这个仓库：

`go get github.com/zzjbattlefield/go_redis`

或者你可以克隆这个仓库到你的本地：

`git clone https://github.com/zzjbattlefield/go_redis.git`

## 使用

你可以使用go命令来运行这个仓库：

`go run main.go`

或者你可以编译成可执行文件：

`go build -o go_redis main.go`

然后运行：

`./go_redis`

你可以使用任何Redis客户端来连接这个服务器，例如redis-cli。默认的端口是6379，你可以在redis.conf文件中修改它。

## 功能

目前，这个仓库实现了以下功能：

- RESP协议的解析和构造
- TCP服务器的创建和监听
- 字典数据结构的实现和操作
- 数据库的创建和切换
- 字符串类型的命令，如SET, GET, DEL等
- AOF持久化的开启和关闭

## 计划

未来，这个仓库计划实现以下功能：

- 其他数据类型的命令，如列表，集合，哈希等
- 事务和管道的支持
- 集群和哨兵的支持
- RDB持久化的支持
