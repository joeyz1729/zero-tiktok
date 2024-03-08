# 项目结构

```shell
tiktok
├─ README.md
├─ apps
│    ├─ app
│    │    └─ tiktok-api
│    ├─ comment
│    │    └─ tiktok-user
│    ├─ favorite
│    │    └─ tiktok-user
│    ├─ follow
│    │    ├─ admin
│    │    └─ tiktok-user
│    ├─ message
│    │    ├─ etc
│    │    ├─ internal
│    │    ├─ message
│    │    ├─ message.go
│    │    ├─ message.proto
│    │    └─ repository
│    ├─ tiktok-user
│    │    ├─ admin
│    │    ├─ repository
│    │    └─ tiktok-user
│    └─ videoservice
│           ├─ admin
│           └─ tiktok-user
├─ doc
│    ├─ first.proto
│    └─ second.proto
├─ go.mod
├─ go.sum
├─ pkg
│    ├─ jwtx
│    │    └─ jwt.go
│    ├─ snowflake
│    │    └─ snowflake.go
│    └─ tool
│           └─ encryption.go
└─ sql
```

## 微服务模块

项目包括用户模块，视频模块，关注模块，点赞模块，评论模块和消息模块。


### API
| server name | port |
|-------------|------|
| api bff     | 8888 |


### RPC
| server name | port |
|-------------|------|
| user        | 9091 |
| video       | 9092 |
| follow      | 9093 |
| favorite    | 9094 |
| comment     | 9095 |
| message     | 9096 |




BFF层负责RESTful API，

微服务板块有user，video，comment，favorite，follow 

BFF层需要先写好api文件，其中包含service内容，request， response，
以及其中使用的结构，例如user，order的结构体

微服务板块中rpc的proto文件也需要先写好对应的request，response和结构体
还有rpc调用函数

BFF logic层处理请求参数并调用rpc函数
rpc函数负责和数据库交互


# BFF
通过api.api编写http路由和handler，然后利用`goctl api go`生成代码。

在svc中添加一个`UserRpc`的client端，用于调用user模块的rpc服务，后续可以添加`VideoRpc, CommentRpc, LikeRpc, FollowRpc`等微服务模块。

修改config和api.yaml添加`RpcClient`的配置，go-zero默认使用etcd作为服务注册和发现，这里修改client配置添加consul。

在程序入口api.go中添加匿名导入`import _ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"`这样在启动后就可以通过consul发现user模块的rpc服务。

# MicroService
以user模块为例，首先编写`user.proto`文件添加需要调用的函数，通过`goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=.`生成go代码。

user模块需要通过consul服务注册，与redis和mysql交互，在生成user id和存储password时使用snowflake和加密，在`config.go`和`user.yaml`中添加对应内容。

在svc中添加mysql和redis实例。

修改`user.go`，在启动服务时添加consul注册和snowflake初始化。



# 
## 基本环境
| Name	   | Description	                   |
|---------|--------------------------------|
| go-zero | web & rpc go frame             |
| Consul  | service registry and discovery |
| Mysql   | database                       |
| Redis   | cache                          |
| Kafka   | message queue                  |
| Docker  | code runtime environment       |


## 启动
1. 启动user模块 `go run user.go`，监听在9091端口，类似启动其他模块。
3. 启动bff层`go run api.go`，监听8888端口。
4. 通过postman测试路由 http://127.0.0.1:8888/douyin/user/register 是否能够正常返回

## 接口介绍

### User

#### Register

#### Login

#### UserInfo

### Video

#### Publish

#### Feed

#### PublishList

获取指定用户发布的视频列表

流程
1. 在bff层首先接收登陆用户的UserId，以及查询用户的AuthorId
2. 通过UserRPC验证AuthorId的正确性，如果正确则返回AuthorInfo
3. 通过FollowRPC查询UserId和AuthorId的关注关系
4. 通过VideoRPC查询AuthorId发布的所有视频Videos详细信息
5. 通过FavoriteRPC查询UserId和VideoIds查询点赞情况
### Favorite

#### FavoriteList

获取用户点赞的视频列表

流程
1. 在bff层首先接收登陆用户的UserId
2. 通过FollowRPC查询UserId所有点赞的视频VideoIds
3. 通过VideoRPC查询VideoIds对应的视频Videos详细信息
4. 通过UserRpc查询AuthorIds的详细信息
5. 通过FollowRpc查询AuthorIds与UserId的关注关系


#### Publish


bff --> user, video --> bff

# 架构

主要思路：mysql用作数据持久，常用数据基本不从mysql中获取。优先从redis-es中。
不能只使用es，还是要加上redis层，必要时还可以加local cache。

## 1. 用户

### 创建
用户成功注册后user表创建条目，通过canal+mq操作user_count表创建，及es中存储用户的全部信息（包含计数）。

### 更新
基础信息目前不涉及更新，计数信息的更新会通过canal+mq更新es。

### 查询
查询用户详细信息，包括基础信息与计数信息。通过redis缓存+es实现。

## 2. 视频

### 创建
通过PublishAction接口，MySQL中video表插入新条目，并通过canal+mq写入video_count表以及es中。

### 更新
点赞（favor）及评论（commit）服务中，更新video_count以及es。
因此需要监听的表是点赞与服务表。

### 查询
视频详细信息通过redis+es查询，不使用mysql。
包括根据发布日期排序查询，根据作者名称查询，根据id/ids指定查询单个/列表。


## 3. 关注

### 创建/删除
关注与取关时，直接操作relation关系表，通过canal更新计数表和es中user的详细信息。
目前没有使用isDel软删除，之后可以考虑优化。

### 更新
关注数与粉丝数的计数需要根据canal+mq更新。

### 查询
查询通过redis+es，不涉及mysql。包括查询关注关系，关注/粉丝列表，计数信息。
关注关系与列表通过redis的set进行实现，不添加超时时间。
计数信息通过es中保存user详细信息实现。


## 4. 点赞

类似


### 5. 评论

类似
