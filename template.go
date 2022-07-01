package main

import (
	"bytes"
	"strings"
	"text/template"
)

var httpTemplate = `
{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
type {{.ServiceType}}EventServer interface {
{{- range .MethodSets}}
	{{.Name}}(context.Context, *{{.Request}}) error
{{- end}}
}

func Register{{.ServiceType}}EventServer(r *message.Router, sub message.Subscriber, srv {{.ServiceType}}EventServer) {
	{{- range .Methods}}
	r.AddNoPublisherHandler(
		"{{.Path}}",
		"{{.Path}}",
		sub,
		_{{$svrType}}_{{.Name}}{{.Num}}_Event_Handler(srv),
	)
	{{- end}}
}

{{range .Methods}}
func _{{$svrType}}_{{.Name}}{{.Num}}_Event_Handler(srv {{$svrType}}EventServer) func(msg *message.Message) error {
	return func(msg *message.Message) error {
		var req {{.Request}}
		err := protojson.Unmarshal(msg.Payload, &req)
		if err != nil {
			return err
		}
		return srv.{{.Name}}(msg.Context(), &req)
	}
}
{{end}}

type {{.ServiceType}}EventClient interface {
{{- range .MethodSets}}
	{{.Name}}(ctx context.Context, req *{{.Request}}) error
{{- end}}
}
	
type {{.ServiceType}}EventClientImpl struct{
	publisher message.Publisher
}
	
func New{{.ServiceType}}EventClient (publisher message.Publisher) {{.ServiceType}}EventClient {
	return &{{.ServiceType}}EventClientImpl{publisher}
}

{{range .MethodSets}}
func (c *{{$svrType}}EventClientImpl) {{.Name}}(ctx context.Context, req *{{.Request}}) error {
	topic := "{{.Path}}"
	byteData, err := protojson.Marshal(req)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), byteData)
	msg.SetContext(ctx)

	return c.publisher.Publish(topic, msg)
}
{{end}}
`

type serviceDesc struct {
	ServiceType string // Greeter
	ServiceName string // helloworld.Greeter
	Metadata    string // api/helloworld/helloworld.proto
	Methods     []*methodDesc
	MethodSets  map[string]*methodDesc
}

type methodDesc struct {
	// method
	Name    string
	Num     int
	Request string
	Reply   string
	// http_rule
	Path         string
	Method       string
	HasVars      bool
	HasBody      bool
	Body         string
	ResponseBody string
}

func (s *serviceDesc) execute() string {
	s.MethodSets = make(map[string]*methodDesc)
	for _, m := range s.Methods {
		s.MethodSets[m.Name] = m
	}
	buf := new(bytes.Buffer)
	tmpl, err := template.New("http").Parse(strings.TrimSpace(httpTemplate))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return strings.Trim(buf.String(), "\r\n")
}
