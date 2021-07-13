module github.com/ii64/protoc-gen-gohttpclient/example

go 1.16

require (
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/ii64/protoc-gen-gohttpclient v0.0.0-00010101000000-000000000000
	github.com/mwitkow/go-proto-validators v0.3.2
	github.com/pkg/errors v0.9.1
	google.golang.org/genproto v0.0.0-20210708141623-e76da96a951f
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.39.0

replace github.com/ii64/protoc-gen-gohttpclient => ../
