package lib

import (
	"net/http"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
)

type ProtoMessageIface interface {
	protoreflect.ProtoMessage
	// proto.Message
	protoiface.MessageV1
}
type ProtoMessageValidator interface {
	Validate() error
}
type HTTPResponseValidatorHandler func(*http.Response) error
type HTTPClientMethodValidatorHandler func(ProtoMessageValidator) error

type HTTPClientService struct {
	BaseURL               string
	Client                *http.Client
	PbDiscardUnknown      bool
	HttpResponseValidator HTTPResponseValidatorHandler
	ResponseValidator     HTTPClientMethodValidatorHandler
}

type HTTPServiceConstructorArg func(*HTTPClientService)

func DefaultHTTPResponseValidator(res *http.Response) error {
	return nil
}
func DefaultMethodValidator(msg ProtoMessageValidator) error {
	return msg.Validate()
}
