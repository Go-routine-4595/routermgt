# Define the target binary name
BINARY_NAME=routermgt

# Define the protoc compiler and flags
# PROTOC := protoc
# PROTOC_FLAGS := --go_opt=paths=source_relative --go_out=. --go-grpc_out=./proto

# Define the source files and proto files
SRCS := $(wildcard *.go)
# PROTO_FILES := $(wildcard proto/*.proto)

# Define the default target
default: build

# proto:
# 	$(PROTOC) $(PROTOC_FLAGS) $(PROTO_FILES)

build: proto
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin $(SRCS)
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux $(SRCS)

run: build
	./${BINARY_NAME}

test:
	go test -v ./...

clean:
	@go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux

.PHONY: build proto clean test run