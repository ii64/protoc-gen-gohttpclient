package lib

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	ErrMethodHasNoHTTPClientSupport  = errors.New("no google.api.http option for this method")
	ErrMethodHasNoHTTPBindingSupport = errors.New("no binding specified for google.api.HttpRule")
)

func (cs *HTTPClientService) ResponseHTTPClientHandler(res *http.Response, dst ProtoMessageIface) (err error) {
	var body []byte
	if body, err = ioutil.ReadAll(res.Body); err != nil {
		return
	}
	var ct string
	if ct, _, err = mime.ParseMediaType(res.Header.Get("Content-Type")); err != nil {
		return
	}
	switch ct {
	case "application/protobuf", "application/x-protobuf":
		if err = proto.Unmarshal(body, dst); err != nil {
			return
		}
	case "application/json", "application/vnd.api+json":
		pj := protojson.UnmarshalOptions{
			DiscardUnknown: cs.PbDiscardUnknown,
		}
		if err = pj.Unmarshal(body, dst); err != nil {
			return
		}
	default:
		return fmt.Errorf("unknwon response content type %q", ct)
	}
	return
}
