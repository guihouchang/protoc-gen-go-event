TEST_PROTO_FILES=$(shell find test/pb -name *.proto)

.PHONY: event-test
# generate event proto
event-test:
	go build -o protoc-gen-go-event && protoc --plugin=./protoc-gen-go-event \
		   --proto_path=./ \
		   --proto_path=./pb \
	       --proto_path=./test/pb \
	       --go_out=paths=source_relative:./ \
	       --go-event_out=paths=source_relative:./ \
	       test/pb/test.proto
