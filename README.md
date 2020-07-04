# grpc-test-demo 文档说明


## 技术说明

go + submodule(共享 git proto 文件库) + grpc

## *.pb.go 文件编译语句

protoc -I. --go_out=plugins=grpc,paths=source_relative:./gen/go/ your/service/v1/your_service.proto

在 go-grpc-proto/prod 终端下执行以下语句生成 prod.pb.go 文件, 其中 go-grpc-proto 是以 submodule 的形式共用第三方 git 库,
prod.proto 引用了 status/status.proto 文件, 故需要 --proto_path=../status/status.proto 编译参数.

protoc -I=. -I=../ --proto_path=../status/status.proto --go_out=plugins=grpc,paths=source_relative:../../src/prod/  prod.proto

在 go-grpc-proto/status 终端下执行以下语句生成 prod.pb.go 文件

protoc -I=.  --go_out=plugins=grpc,paths=source_relative:../../src/status status.proto

## 使用步骤

### go 到 go

git clone https://github.com/tanjunchen/grpc-test-demo.git
    
git submodule update --init --recursive
