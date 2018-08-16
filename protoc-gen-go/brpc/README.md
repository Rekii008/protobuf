
给proto插上rpc的翅膀，向grpc致敬！



目前暂时只提供http服务，底层的驱动只有gin

计划加上真正的rpc（包括client和server），底层驱动使用net/http，欢迎pr



使用方式：

1. 定义proto

   ```
   syntax = "proto3";

   import "github.com/golang/protobuf/protoc-gen-go/brpc/brpc.proto";  //为method扩展了选项

   package helloworld;

   service Greeter {   //可以多个service XX 定义，每个对应生成一个 ServerXX 服务
       rpc SayHello (HelloRequest) returns (HelloReply) {  //ServerXX 服务需要注册实现的接口
           option(brpc.Method) = "POST";    //接口的HTTP方法, POST, GET, OPTION等
           option(brpc.Path) = "/greeter/sayhello";   //接口的HTTP路径
       }
   }

   message HelloRequest {
       string name = 1;
   }

   message HelloReply {
       string message = 1;
   }
   ```

2. 实现proto中，服务定义的所有接口

   ```
   type Example struct {
   }

   func (e *Example) SayHello(c brpc.Context, req *helloworld.HelloRequest) (rsp *helloworld.HelloReply, err brpc.Error) {
       if req.Name == "rek" {
           return &helloworld.HelloReply{Message: "hello, " + req.Name}, nil
       }
       return nil, brpc.NewError(100, "invalid req name")
   }
   ```

3. 生成服务，注册服务

   ```
   e := helloworld.NewGreeterServer()
   e.Register(&Example{})
   ```

4. 启动服务

   ```
   e.Run(":8080")
   ```

5. 也可以注册其他服务，这是未来实现真正rpc的基础

   ```
   e.Engine.Register("GET", "/other", Other)
   func Other(c brpc.Context) {
       c.Status(200)
       c.ResponseString("other is ok")

   }
   ```

以上代码来自Example