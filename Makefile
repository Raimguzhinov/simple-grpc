PHONY: generate

generate:
	mkdir -p pkg
	protoc --proto_path=api --go_out=. \
		--go_opt=module=github.com/Raimguzhinov/simple-grpc \
		--go-grpc_out=. --go-grpc_opt=module=github.com/Raimguzhinov/simple-grpc \
		api/protobuf/eventmanager.proto
	
server:
	go build -o ./build/event_server ./cmd/server
	./build/event_server -h 127.0.0.1 -p 50051

client: 
	go build -o ./build/event_client ./cmd/client
	./build/event_client -dst 127.0.0.1 -p 50051 -sender-id 400

broker:
	docker run -it --rm --detach --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management
