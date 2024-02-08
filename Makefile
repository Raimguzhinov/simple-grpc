PHONY: generate

generate:
	mkdir -p pkg
	protoc --go_out=pkg --go_opt=paths=source_relative \
		--go-grpc_out=pkg --go-grpc_opt=paths=source_relative \
		api/protobuf/eventmanager.proto

server:
	go run ./cmd/server -h 127.0.0.1 -p 50051

client: 
	go run ./cmd/client -dst 127.0.0.1 -p 50051 -sender-id 400
