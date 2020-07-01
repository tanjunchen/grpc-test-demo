# protoc 文档

protoc -I. --go_out=plugins=grpc,paths=source_relative:./gen/go/ your/service/v1/your_service.proto


protoc --proto_path=IMPORT_PATH1  --proto_path=IMPORT_PATH2 --cpp_out=DST_DIR --java_out=DST_DIR --python_out=DST_DIR 
--go_out=DST_DIR --ruby_out=DST_DIR --objc_out=DST_DIR --csharp_out=DST_DIR path/to/file.proto

https://www.cnblogs.com/FireworksEasyCool/p/12782137.html

# 测试全量生成 proto 文件

protoc --go_out=plugins=grpc:. grpc-test-demo/go-grpc-proto/**/*.proto

# 问题

go-grpc-proto 是一个 submodule

1. 怎么在 grpc-go-test 项目根目录下将 grpc-go-test 下的 proto 文件生成 pb 文件到 src/proto/ 文件下

1. 怎么在 go-grpc-proto/example/example.proto 中引用 go-grpc-proto/status/status.proto 中的 proto 文件

1. 怎么在 go-grpc-proto/example/example.proto 中引用 google/api/http.proto 文件

