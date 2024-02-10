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