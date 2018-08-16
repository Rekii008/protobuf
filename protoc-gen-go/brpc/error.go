package brpc

import "github.com/golang/protobuf/proto"

type Error interface {
    String() string
}

func NewError(errcode int32, err interface{}) Error {
    switch t := err.(type) {
    case string:
        return &SysError{Errcode: proto.Int32(errcode), Errmsg: proto.String(t)}
    case error:
        return &SysError{Errcode: proto.Int32(errcode), Errmsg: proto.String(t.Error())}
    case Error:
        return &SysError{Errcode: proto.Int32(errcode), Errmsg: proto.String(t.String())}
    default:
        return &SysError{Errcode: proto.Int32(errcode), Errmsg: proto.String("unknown error")}
    }
}

const (
    WriteError      = int32(-100)
    UnmarshalError  = int32(-101)
    MarshalError    = int32(-102)
)
