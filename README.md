# grpc-test-demo 文档说明

通过 submodule 异构语言共用一套 proto 文件, 在 proto 中 import 其他的 proto 文件.

## 技术说明

go + submodule(共享 git proto 文件库) + grpc

## *.pb.go 文件编译语句

protoc -I. --go_out=plugins=grpc,paths=source_relative:./gen/go/ your/service/v1/your_service.proto

在 go-grpc-proto/prod 终端下执行以下语句生成 prod.pb.go 文件, 其中 go-grpc-proto 是以 submodule 的形式共用第三方 git 库,
prod.proto 引用了 status/status.proto 文件, 故需要 --proto_path=../status/status.proto 编译参数.

protoc -I=. -I=../ --proto_path=../status/status.proto --go_out=plugins=grpc,paths=source_relative:../../src/prod/  prod.proto

在 go-grpc-proto/status 终端下执行以下语句生成 prod.pb.go 文件

protoc -I=.  --go_out=plugins=grpc,paths=source_relative:../../src/status status.proto

java 可以直接执行 proto-controller 下的 protobuf:compile 与 protobuf:compile-custom 命令.

## 目录结构

```
├── client      # go 测试客户端
│   └── main.go
├── go-grpc-proto       # 存放原始 proto 文件
│   ├── prod
│   │   └── prod.proto      # 说明 prod.proto 引用了 status 下的包
│   └── status
│       └── status.proto
├── go.mod      # mod 文件
├── README.md       # readme 文件
├── server      # go 服务器后端
│   └── main.go
├── service     # demo 示例
│   └── test_service.go
└── src     # 生成后的 *.pb.go 文件
    ├── prod
    │   └── prod.pb.go
    └── status
        └── status.pb.go
```

## 使用步骤

### go 到 go

git clone https://github.com/tanjunchen/grpc-test-demo.git
    
git submodule update --init --recursive

go mod tidy

cd server && go run main.go

```
Listen on 0.0.0.0:9999
```

server 端

```
Listen on 0.0.0.0:9999
211
```

client 端

```
prod_stock:211  status:{code:"200"  msg:"success"}
```

### java 到 go

参见 java [tanjunchen/Java-Go-Grpc-Demo](https://github.com/tanjunchen/Java-Go-Grpc-Demo)

启动服务端同理
 
git clone https://github.com/tanjunchen/Java-Go-Grpc-Demo.git
 
git submodule update --init --recursive
 
cd src/main/java/com/ctj/SpringBootGRPCApplication.java

终端响应如下：

```
xxxxxx  INFO 20752 --- [           main] o.a.c.c.C.[Tomcat].[localhost].[/]       : Initializing Spring embedded WebApplicationContext
xxxxxx  INFO 20752 --- [           main] o.s.web.context.ContextLoader            : Root WebApplicationContext: initialization completed in 1696 ms
xxxxxx  INFO 20752 --- [           main] o.s.s.concurrent.ThreadPoolTaskExecutor  : Initializing ExecutorService 'applicationTaskExecutor'
xxxxxx  INFO 20752 --- [           main] o.s.b.w.embedded.tomcat.TomcatWebServer  : Tomcat started on port(s): 9998 (http) with context path ''
xxxxxx  INFO 20752 --- [           main] com.ctj.SpringBootGRPCApplication        : Started SpringBootGRPCApplication in 2.513 seconds (JVM running for 2.879)
prod_stock: 1000
status {
  code: "200"
  msg: "success"
}

200
success
```

java 推荐使用 https://github.com/yidongnan/grpc-spring-boot-starter 与 grpc 服务提供者通信.