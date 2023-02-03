package encoder

import (
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
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

// responseEncoder encodes the object to the HTTP response.
func ResponseEncoder(w http.ResponseWriter, r *http.Request, data interface{}) error {
	type Response struct {
		Code    int         `json:"code"`
		Results interface{} `json:"results"`
		Message string      `json:"message"`
	}
	res := &Response{
		Code:    200,
		Results: data,
		Message: "success",
	}
	msRes, err := json.Marshal(&res)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(se.Code)
	_, _ = w.Write(msRes)
	return nil
}

// ContentType returns the content-type with base prefix.
func ContentType(subtype string) string {
	return strings.Join([]string{baseContentType, subtype}, "/")
}
