PHONY: generate

generate:
	mkdir -p pkg
	protoc --go_out=pkg --go_opt=paths=source_relative \
		--go-grpc_out=pkg --go-grpc_opt=paths=source_relative \
		api/protobuf/eventmanager.proto
