# GRPC Golang 范例

> 参考文档 <https://grpc.io/docs/languages/go/quickstart/>
>
> 以下内容的工作目录始终不变

在 `./pb/` 下创建 greeter.proto 文件，定义 protobuf。

```protobuf
syntax = "proto3";

package grpc_example;

option go_package = ".;grpc_example";

service Greeter {
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}
```

安装 protoc，各 OS 平台安装方法不同。

使用 go install 安装 protoc-gen-go 和 protoc-gen-grpc，并把装好的可执行文件路径 `$GOPATH/bin` 加入 PATH 环境变量。

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
export PATH="$PATH:$(go env GOPATH)/bin"
```

使用上述安装的三个工具生成两个 .go 文件 `greeter.pb.go` 和 `greeter_grpc.pb.go`，不要自行编辑这两个文件。

如果 .proto 文件有变动，则使用下面的命令重新生成 .go 文件

```sh
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pb/greeter.proto
```
