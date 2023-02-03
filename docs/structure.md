# 项目结构

## 调整 kratos 项目结构

创建 kratos 项目后，我们可以看到项目结构如下

```shell
application
|____api
| |____helloworld
| | |____v1
|____cmd
| |____helloworld
|____configs
|____internal
| |____conf
| |____data
| |____biz
| |____service
| |____server
|____test
|____third_party
|____pkg
|____go.mod
|____go.sum
|____LICENSE
|____README.md
```

使用Mono-repo大仓库，项目结构调整为如下结构

```shell
application
|____api
| |____user
| | |____service
| | | |____v1
|____app
| |____user
| | |____admin
| | |____interface
| | |____job
| | |____service
| | | |____cmd
| | | |____configs
| | | |____internal
| | | | |____conf
| | | | |____data
| | | | |____biz
| | | | |____service
| | | | |____server
| | | | |____test
|____third_party
|____third_party
|____pkg
|____go.mod
|____go.sum
|____Makefile
|____LICENSE
|____README.md
```
### 根目录

#### /api
API 协议定义目录，services.proto protobuf 文件，以及生成的 go 文件。我们通常把 api 文档直接在 proto 文件中描述。

#### /pkg
外部应用程序可以使用的库代码,注意，internal 目录是确保私有包不可导入的更好方法，因为它是由 Go 强制执行的
### 应用类型目录

kratos 把微服务中的 app 服务类型主要分为5类：interface、service、job、admin、task，，应用 cmd 目录负责程序的：启动、关闭、配置初始化等。

app/user/下面的一级目录就是应用类型目录

1. interface: 对外的 BFF 服务，接受来自用户的请求，比如暴露了 HTTP/gRPC 接口。
2. service: 对内的微服务，仅接受来自内部其他服务或者网关的请求，比如暴露了gRPC 接口只对内服务。
3. admin：区别于 service，更多是面向运营测的服务，通常数据权限更高，隔离带来更好的代码级别安全。
4. job: 流式任务处理的服务，上游一般依赖 message broker。
5. task: 定时任务，类似 cronjob，部署到 task 托管平台中。


## Mono repo开发流程
1. 往mono repo添加一个服务
```shell
// 新建一个kratos服务
kratos new app/demo/service --nomod

// 我们也可以将生成的服务指定到合适的服务类型目录下面，比如
kratos new app/demo/interface --nomod
```

2. 生成编译服务api proto文件
```shell
// 生成
kratos proto add api/demo/service/v1/demo.proto 
// 编译
kratos proto client api/demo/service/v1/demo.proto
```

3. 生成service代码（也就是mvc中的controller）
```shell
kratos proto server api/demo/service/v1/demo.proto -t app/demo/service/internal/service
```

# 参考文档:
[Go工程化 - Project Layout 最佳实践](https://go-kratos.dev/blog/go-project-layout)