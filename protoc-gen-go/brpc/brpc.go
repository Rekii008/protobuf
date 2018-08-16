package brpc

import (
    "github.com/golang/protobuf/protoc-gen-go/generator"
    pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
    "fmt"
    "path"
    "github.com/golang/protobuf/proto"
)

func init() {
    generator.RegisterPlugin(&brpc{})
}

type brpc struct {
    gen *generator.Generator
}

func (e *brpc) p(args ...interface{}) {
    e.gen.P(args...)
}

func (e *brpc) e(msg ...string) {
    e.gen.Fail(msg...)
}

func (e *brpc) objectNamed(name string) generator.Object {
    e.gen.RecordTypeUse(name)
    return e.gen.ObjectNamed(name)
}

func (e *brpc) typeName(str string) string {
    return e.gen.TypeName(e.objectNamed(str))
}

func (e *brpc) Name() string {
    return "brpc"
}

func (e *brpc) Init(g *generator.Generator) {
    e.gen = g
}

func (e *brpc) GenerateImports(file *generator.FileDescriptor) {
    if len(file.FileDescriptorProto.Service) == 0 {
        return
    }
    e.p("// Import for brpc")
    e.p("import brpc", " ", generator.GoImportPath(path.Join(string(e.gen.ImportPrefix), "github.com/golang/protobuf/protoc-gen-go/brpc")))
}

func (e *brpc) Generate(file *generator.FileDescriptor) {
    if len(file.FileDescriptorProto.Service) == 0 {
        return
    }
    e.p("// Services of brpc")
    for _, service := range file.FileDescriptorProto.Service {
        e.generateService(file, service)
    }
}

type method struct {
    cmdid uint32
    name string
    input string
    output string
    method string
    path string
}

func (e *brpc) generateService(file *generator.FileDescriptor, service *pb.ServiceDescriptorProto) {
    serviceName := service.GetName()
    serverName := serviceName + "Server"
    serviceInterface := "I" + serviceName
    methods := make([]method, 0)
    for _, mm := range service.Method {
        m := method{
            cmdid: 0,
            name: mm.GetName(),
            input: e.typeName(mm.GetInputType()),
            output: e.typeName(mm.GetOutputType()),
        }
        option := mm.GetOptions()
        if method, err := proto.GetExtension(option, E_Method); err != nil {
            m.method = "POST"
        } else {
            m.method = *method.(*string)
        }
        if path, err := proto.GetExtension(option, E_Path); err != nil {
            e.e("empty path of ", mm.GetName())
        } else {
            m.path = *path.(*string)
        }
        methods = append(methods, m)
    }

    e.p("// Server of the " + serviceName + "Service")
    e.p("type " + serviceInterface + " interface {")
    for _, m := range methods {
        e.p("// register on: " + m.method + " " + m.path)
        e.p(fmt.Sprintf(`%s (c brpc.Context, req * %s) (rsp * %s, err brpc.Error)`, m.name, m.input, m.output))
    }
    e.p("}")
    e.p()
    e.p("type " + serverName + " struct{")
    e.p("brpc.Engine")
    e.p("}")
    e.p()
    e.p("func New" + serverName + "() *" + serverName + "{")
    e.p("return &" + serverName + "{")
    e.p("Engine: brpc.NewEngine(),")
    e.p("}")
    e.p("}")
    e.p()
    e.p("func (s *" + serverName + ") Run(address ...string) brpc.Error {")
    e.p("return s.Engine.Run(address...)")
    e.p("}")
    e.p()
    e.p("func (s *" + serverName + ") Register(i " + serviceInterface + ") {")
    for _, m := range methods {
        e.generateRegisterMethod(&m)
    }
    e.p("}")

}

func (e *brpc) generateRegisterMethod(m *method) {
    e.p(fmt.Sprintf(`
s.Engine.Register("%s", "%s", func (c brpc.Context) {
    var req %s
		if e := c.GetProto(&req); e != nil {
			c.ResponseError(e)
            return
		}

		rsp, e := i.%s(c, &req)
		if e != nil {
			c.ResponseError(e)
            return
		}

		if e := c.ResponseProto(rsp); e != nil {
            c.ResponseError(e)
            return
        }
})`, m.method, m.path, m.input, m.name))
}















