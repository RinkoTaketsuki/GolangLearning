# gRPC 用法总结

## `.pb.go` 代码生成

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

protoc --go_out=out/path/ --go_opt=paths=source_relative \
    --go-grpc_out=out/path/ --go-grpc_opt=paths=source_relative \
    path/to/file/idl.proto
```

## 指定 go package

```proto
option go_package = "my/go/mod/path/to/pkg";
```

## 初始化 client

```go
	var opts []grpc.DialOption
	// initialize opts ...
	conn, err := grpc.NewClient(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewRouteGuideClient(conn)
```

## 初始化 server

```go
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	// initialize opts ...
	grpcServer := grpc.NewServer(opts...)
	// newServer() initailize an object implementing pb.ServiceNameServer interface
	pb.RegisterRouteGuideServer(grpcServer, newServer())
	grpcServer.Serve(lis)
```

## 4 类 RPC 的用法

### simple

直接调用即可。

### server-side streaming

* client

无限循环 `stream.Recv()` 接收数据直到 `err` 为 `io.EOF`。

* server

使用 `stream.Send()` 逐个发送数据，整个函数返回时结束发送。

### client-side streaming

* client

使用 `stream.Send()` 逐个发送数据，用 `stream.CloseAndRecv()` 结束发送并接收返回数据。

* server

无限循环 `stream.Recv()` 接收数据直到 `err` 为 `io.EOF`。最后用 `stream.SendAndClose()` 发送返回数据并结束接收。

### bidirectional streaming

* client

发送流使用 `stream.Send()` 逐个发送数据，用 `stream.CloseSend()` 结束发送。
接收流无限循环 `stream.Recv()` 接收数据直到 `err` 为 `io.EOF`。

* server

发送流使用 `stream.Send()` 逐个发送数据，整个函数返回时结束发送。
接收流无限循环 `stream.Recv()` 接收数据直到 `err` 为 `io.EOF`。
