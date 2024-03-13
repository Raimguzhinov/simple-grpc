.PHONY: generate
generate:
	mkdir -p pkg
	protoc --proto_path=api --go_out=. \
		--go_opt=module=github.com/Raimguzhinov/simple-grpc \
		--go-grpc_out=. --go-grpc_opt=module=github.com/Raimguzhinov/simple-grpc \
		api/protobuf/eventmanager.proto

.PHONY: compile
compile:
	go build -o ./build/ ./...

.PHONY: run
run:
	./build/server

.PHONY: broker
broker:
	docker start rabbitmq || docker run -it --rm --detach --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management

.PHONY: test
test:
	gotestsum --format pkgname --raw-command go test -json -cover ./...
