# 基于protobuff生成事件api接口与客户端
## 定义事件
```protobuf
syntax = "proto3";

package test;

import "event/options/event.proto";


option go_package = "code.aliyun.com/fz.7799520.com/protoc-gen-go-event/test/pb";

service TestService {
  rpc TestMethod (TestRequest) returns (TestResponse) {
    option (event.eventName) = "/test/test";
  }
}

message TestRequest {

}

message TestResponse {

}
```
# 生成接口
```shell
   protoc  --proto_path=./ \
		   --proto_path=./pb \
	       --proto_path=./test/pb \
	       --go_out=paths=source_relative:./ \
	       --go-event_out=paths=source_relative:./ \
	       test/pb/test.proto
```
# 生成接口文件
```go
// Code generated by protoc-gen-go-event. DO NOT EDIT.
// versions:
// protoc-gen-go-event v2.3.1

package pb

import (
	context "context"
	watermill "github.com/ThreeDotsLabs/watermill"
	amqp1 "github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	message "github.com/ThreeDotsLabs/watermill/message"
	amqp "github.com/streadway/amqp"
	protojson "google.golang.org/protobuf/encoding/protojson"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = new(amqp.Table)
var _ = new(amqp1.Config)
var _ = new(message.Message)
var _ = watermill.NewUUID
var _ = protojson.Marshal

type TestServiceEventServer interface {
	TestMethod(context.Context, *TestRequest) error
	TestMethod1(context.Context, *TestRequest) error
}

func RegisterTestServiceEventServer(r *message.Router, sub message.Subscriber, srv TestServiceEventServer) {
	r.AddNoPublisherHandler(
		"",
		"",
		sub,
		_TestService_TestMethod10_Event_Handler(srv),
	)
	r.AddNoPublisherHandler(
		"/test/test",
		"/test/test",
		sub,
		_TestService_TestMethod0_Event_Handler(srv),
	)
}

func _TestService_TestMethod10_Event_Handler(srv TestServiceEventServer) func(msg *message.Message) error {
	return func(msg *message.Message) error {
		var req TestRequest
		err := protojson.Unmarshal(msg.Payload, &req)
		if err != nil {
			return err
		}
		return srv.TestMethod1(msg.Context(), &req)
	}
}

func _TestService_TestMethod0_Event_Handler(srv TestServiceEventServer) func(msg *message.Message) error {
	return func(msg *message.Message) error {
		var req TestRequest
		err := protojson.Unmarshal(msg.Payload, &req)
		if err != nil {
			return err
		}
		return srv.TestMethod(msg.Context(), &req)
	}
}

type TestServiceEventClient interface {
	TestMethod(ctx context.Context, req *TestRequest) error
	TestMethod1(ctx context.Context, req *TestRequest) error
}

type TestServiceEventClientImpl struct {
	publisher message.Publisher
}

func NewTestServiceEventClient(publisher message.Publisher) TestServiceEventClient {
	return &TestServiceEventClientImpl{publisher}
}

func (c *TestServiceEventClientImpl) TestMethod(ctx context.Context, req *TestRequest) error {
	topic := "/test/test"
	byteData, err := protojson.Marshal(req)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), byteData)
	msg.SetContext(ctx)

	return c.publisher.Publish(topic, msg)
}

func (c *TestServiceEventClientImpl) TestMethod1(ctx context.Context, req *TestRequest) error {
	topic := ""
	byteData, err := protojson.Marshal(req)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), byteData)
	msg.SetContext(ctx)

	return c.publisher.Publish(topic, msg)
}

```
