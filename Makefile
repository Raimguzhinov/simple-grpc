PHONY: generate

generate:
	mkdir -p pkg
	protoc --go_out=pkg --go_opt=paths=source_relative \
		--go-grpc_out=pkg --go-grpc_opt=paths=source_relative \
		api/protobuf/eventmanager.proto

server:
	go build -o ./build/event_server ./cmd/server
	./build/event_server -h 127.0.0.1 -p 50051

client: 
	go build -o ./build/event_client ./cmd/client
	./build/event_client -dst 127.0.0.1 -p 50051 -sender-id 400
