package brpc

import (
    "time"
    "net/http"
    "github.com/golang/protobuf/proto"
    "io"
)

type Context interface {
    GetString(key string) (s string)
    GetBool(key string) (b bool)
    GetInt(key string) (i int)
    GetInt64(key string) (i64 int64)
    GetFloat64(key string) (f64 float64)
    GetTime(key string) (t time.Time)
    GetDuration(key string) (d time.Duration)
    GetStringSlice(key string) (ss []string)
    GetStringMap(key string) (sm map[string]interface{})
    GetHeader(key string) string
    Cookie(name string) (string, error)

    Param(key string) string
    Query(key string) string
    DefaultQuery(key, defaultValue string) string

    Status(code int)
    Header(key, value string)
    SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)

    GetProto(message proto.Message) Error
    Request() io.Reader
    GetRequest() *http.Request

    ResponseProto(message proto.Message) Error
    ResponseError(err Error) Error
    ResponseString(str string) Error
    ResponseBytes(byt []byte) Error
    Response(reader io.Reader) Error

    GetResponse() http.ResponseWriter
}

type Engine interface {
    Register(method, path string, f func(c Context))
    Run(address ...string) Error
}

