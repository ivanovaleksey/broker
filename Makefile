PATH_PROTO_FAST="pkg/pb/broker_fast"
PATH_PROTO="pkg/pb/broker"

.PHONY: test
test:
	go test -v -count=1 ./...

.PHONY: build
build:
	go build -o ./bin/mbs ./cmd/broker

.PHONY: build-producer
build-producer:
	go build -o ./bin/mbp ./cmd/producer

.PHONY: build-consumer
build-consumer:
	go build -o ./bin/mbc ./cmd/consumer

.PHONY: build-all
build-all: build build-consumer build-producer

.PHONY: generate-proto-fast
generate-proto-fast:
	mkdir -p $(PATH_PROTO_FAST) && cd $(PATH_PROTO_FAST) && \
	protoc \
		-I ../../.. \
		--gogofaster_out=plugins=grpc:. \
		../../../spec.proto

.PHONY: generate-proto
generate-proto:
	mkdir -p $(PATH_PROTO) && cd $(PATH_PROTO) && \
	protoc \
		-I ../../.. \
		--go_out=plugins=grpc:. \
		../../../spec.proto

