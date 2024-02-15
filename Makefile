PHONY: generate

generate:
	mkdir -p pkg
	protoc --proto_path=api --go_out=. \
		--go_opt=module=github.com/Raimguzhinov/simple-grpc \
		--go-grpc_out=. --go-grpc_opt=module=github.com/Raimguzhinov/simple-grpc \
		api/protobuf/eventmanager.proto

compile:
	go build -o ./build/ ./...

run:
	./build/server

broker:
	docker start rabbitmq || docker run -it --rm --detach --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management

test:
	gotestsum --format pkgname --raw-command go test -json -cover ./...
