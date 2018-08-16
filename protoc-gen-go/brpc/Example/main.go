package main

import (
    "github.com/golang/protobuf/protoc-gen-go/brpc"
    "github.com/golang/protobuf/protoc-gen-go/brpc/Example/helloworld"
)

//go:generate mkdir -p helloworld
//go:generate protoc --go_out=plugins=brpc:./helloworld -I$GOPATH/src:. ./example.proto

type Example struct {
}

func (e *Example) SayHello(c brpc.Context, req *helloworld.HelloRequest) (rsp *helloworld.HelloReply, err brpc.Error) {
    if req.Name == "rek" {
        return &helloworld.HelloReply{Message: "hello, " + req.Name}, nil
    }
    return nil, brpc.NewError(100, "invalid req name")
}

func Other(c brpc.Context) {
    c.Status(200)
    c.ResponseString("other is ok")

}

func main() {
    e := helloworld.NewGreeterServer()
    e.Register(&Example{})
    e.Engine.Register("GET", "/other", Other)
    e.Run(":8080")
}