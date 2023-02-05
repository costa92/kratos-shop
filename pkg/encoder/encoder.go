package encoder

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	nethttp "net/http"
	"strings"
)

const (
	baseContentType = "application"
)

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewResponseError(code int, message string) *ResponseError {
	return &ResponseError{
		Code:    code,
		Message: message,
	}
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("ResponseError: %d", e.Code)
}

func FromError(err error) *ResponseError {
	if err != nil {
		return nil
	}
	if se := new(errors.Error); errors.As(err, &se) {
		return NewResponseError(int(se.Code), se.Message)
	}
	return &ResponseError{
		Code:    500,
		Message: err.Error(),
	}
}

func ErrorEncoder(w nethttp.ResponseWriter, r *nethttp.Request, err error) {
	se := errors.FromError(err)
	codec, _ := http.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(se)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/"+codec.Name())
	// w.WriteHeader(se.Code)
	_, _ = w.Write(body)
}

type BaseResponse struct {
	Code int32  `json:"code"`
	Msg  string `json:"message"`
}

// responseEncoder encodes the object to the HTTP response.
func ResponseEncoder(w nethttp.ResponseWriter, r *nethttp.Request, i interface{}) error {
	reply := &Response{
		Code: 200,
	}
	if m, ok := i.(proto.Message); ok {
		payload, err := anypb.New(m)
		if err != nil {
			return err
		}
		reply.Data = payload
	}

	codec := encoding.GetCodec("json")
	data, err := codec.Marshal(reply)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	return nil
}

// ContentType returns the content-type with base prefix.
func contentType(subtype string) string {
	return strings.Join([]string{baseContentType, subtype}, "/")
}
