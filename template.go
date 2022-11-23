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

func Register{{.ServiceType}}EventServer(r *message.Router, sg func(topic string) message.Subscriber, srv {{.ServiceType}}EventServer) {
	{{- range .Methods}}
	r.AddNoPublisherHandler(
		"{{.EventName}}",
		"{{.EventName}}",
		sg("{{.EventName}}"),
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
	{{if gt .EventDelay 0 }}
    {{.Name}}WithDelay(ctx context.Context, req *{{.Request}}, delay int32) error
	{{end}}
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
	topic := "{{.EventName}}"
	byteData, err := protojson.Marshal(req)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), byteData)
	msg.SetContext(ctx)
	{{if gt .EventDelay 0}}
		// 设置延迟队列时间，单位为{{.EventDelay}}ms
		msg.Metadata.Set("x-delay", "{{.EventDelay}}")
	{{end}}
	return c.publisher.Publish(topic, msg)
}
{{end}}
{{range .MethodSets}}
{{if gt .EventDelay 0}}
func (c *{{$svrType}}EventClientImpl) {{.Name}}WithDelay(ctx context.Context, req *{{.Request}}, delay int32) error {
	topic := "{{.EventName}}"
	byteData, err := protojson.Marshal(req)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), byteData)
	msg.SetContext(ctx)
	msg.Metadata.Set("x-delay", fmt.Sprintf("%d", delay))
	return c.publisher.Publish(topic, msg)
}
{{end}}
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
	// event_rule
	EventName  string // 事件名称
	EventDelay int32  // 延迟时间
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
