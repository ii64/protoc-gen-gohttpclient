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

type HTTPRequestPreflightHandler func(req *http.Request) (res *http.Request, err error)
type HTTPResponseValidatorHandler func(res *http.Response) (err error)
type HTTPClientMethodValidatorHandler func(ProtoMessageValidator) (err error)
type HTTPServiceConstructorArg func(*HTTPClientService)

type HTTPClientService struct {
	BaseURL               string
	Client                *http.Client
	PbDiscardUnknown      bool
	HttpRequestPreflight  HTTPRequestPreflightHandler
	HttpResponseValidator HTTPResponseValidatorHandler
	ResponseValidator     HTTPClientMethodValidatorHandler
}

func (c HTTPClientService) Do(req *http.Request) (res *http.Response, err error) {
	req, err = c.HttpRequestPreflight(req)
	if err != nil {
		return
	}
	return c.Client.Do(req)
}

func DefaultHTTPResponseValidator(res *http.Response) error {
	return nil
}
func DefaultMethodValidator(msg ProtoMessageValidator) error {
	return msg.Validate()
}
