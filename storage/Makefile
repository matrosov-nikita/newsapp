build:
	protoc --go_out=plugins=grpc:. proto/storage.proto
	go build -o ../bin/storage ../cmd/storage
