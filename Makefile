.PHONY: build run-server run proto

build:
	go build

run-server: build
	./debugger server \
		--connection-string "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" \
		--kubeconfig $$HOME/.kube/config

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    	   --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/debugger.proto