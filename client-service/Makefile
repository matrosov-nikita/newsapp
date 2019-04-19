build:
	protoc --go_out=plugins=grpc:. proto/newsclient.proto
	go build -o ../bin/client ../cmd/client
