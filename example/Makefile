generate:
	protoc \
		-I . \
		-I ./googleapis/googleapis \
		--go_out=paths=source_relative:./gen/ \
		--go-grpc_out=paths=source_relative:./gen/ \
		--grpc-gateway_out=allow_patch_feature=false,paths=source_relative:./gen/ \
		--gohttpclient_out=paths=source_relative:./gen/ \
		./service.proto