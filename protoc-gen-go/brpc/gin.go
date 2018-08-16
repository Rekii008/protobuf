package brpc

import (
    "github.com/gin-gonic/gin"
    "log"
    "net/http"
    "io"
    "github.com/golang/protobuf/proto"
    "github.com/golang/protobuf/jsonpb"
    "bytes"
)

type ginContext struct {
    *gin.Context
}

func (c *ginContext) GetProto(message proto.Message) Error {
    unmarshaler := jsonpb.Unmarshaler{AllowUnknownFields: true}
    if e := unmarshaler.Unmarshal(c.Request(), message); e != nil {
        return NewError(UnmarshalError, e)
    }
    return nil
}

func (c *ginContext) Request() io.Reader {
    return c.GetRequest().Body
}
func (c *ginContext) GetRequest() *http.Request {
    return c.Context.Request
}

func (c *ginContext) ResponseBytes(byt []byte) Error {
    if _, e := c.GetResponse().Write(byt); e != nil {
        return NewError(WriteError, e)
    }
    return nil
}
func (c *ginContext) ResponseString(str string) Error {
    return c.ResponseBytes([]byte(str))
}
func (c *ginContext) ResponseError(err Error) Error {
    if v, ok := err.(proto.Message); ok {
        return c.ResponseProto(v)
    }
    return c.ResponseString(err.String())
}
func (c *ginContext) Response(reader io.Reader) Error {
    buf := &bytes.Buffer{}
    buf.ReadFrom(reader)
    return c.ResponseBytes(buf.Bytes())
}
func (c *ginContext) ResponseProto(message proto.Message) Error {
    marshaler := jsonpb.Marshaler{
        EmitDefaults: true,
        OrigName:     true,
    }
    buf := &bytes.Buffer{}
    if e := marshaler.Marshal(buf, message); e != nil {
        return NewError(MarshalError, e)
    }
    return c.ResponseBytes(buf.Bytes())
}
func (c *ginContext) GetResponse() http.ResponseWriter {
    return c.Context.Writer
}



type ginEngine struct {
    *gin.Engine
}

func (e *ginEngine) Register(method, path string, f func(c Context)) {
    ginFunc := func (gc *gin.Context) {
        f(&ginContext{Context: gc})
    }
    if method == "POST" {
        e.POST(path, ginFunc)
    } else if method == "GET" {
        e.GET(path, ginFunc)
    } else {
        log.Fatal("method not supported", method)
    }
}

func (e *ginEngine) Run(address ...string) Error {
    if err := e.Engine.Run(address...); err != nil {
        return NewError(-1, "run fail")
    }
    return nil
}

func NewEngine() Engine {
    return &ginEngine{gin.Default()}
}