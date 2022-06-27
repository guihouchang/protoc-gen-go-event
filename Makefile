TEST_PROTO_FILES=$(shell find test/pb -name *.proto)

.PHONY: event
# generate event proto
event:
	cd ./cmd && go build -o protoc-gen-go-event && cd ../ && protoc --plugin=./cmd/protoc-gen-go-event \
		   --proto_path=./ \
		   --proto_path=./pb \
	       --proto_path=./test/pb \
	       --go_out=paths=source_relative:./ \
	       --go-event_out=paths=source_relative:./ \
	       test/pb/test.proto