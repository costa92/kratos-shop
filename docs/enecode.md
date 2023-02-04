# HTTP错误返回结构自定义

##  kratos 默认返回结构体

```json
{
  "id": "0",
  "name": "11"
}
```

这个结构在接口返回是无法满足的，需要修改返回如下格式

```json
{
  "code": 200,
  "results": {
    "name": "11"
  },
  "message": "success"
}
```

## 修改方法
在 kratos 的 transport/http 包中，为我们提供了 http 的相关方法，我们今天就是通过 transport/http 包中的 http.ResponseEncoder 控制,该函数需要传入 EncodeResponseFunc

在 transport/http/coder.go 文件中有下面两个方法: **DefaultResponseEncoder** 与 **DefaultErrorEncoder**
**DefaultResponseEncoder**: 是否默认正常的返回格式
**DefaultErrorEncoder**: 是否默认错误的返回格式

```go
type DecodeRequestFunc func(*http.Request, interface{}) error

// EncodeResponseFunc is encode response func.
type EncodeResponseFunc func(http.ResponseWriter, *http.Request, interface{}) error

// EncodeErrorFunc is encode error func.
type EncodeErrorFunc func(http.ResponseWriter, *http.Request, error)


// DefaultResponseEncoder encodes the object to the HTTP response.
func DefaultResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if v == nil {
		return nil
	}
	if rd, ok := v.(Redirector); ok {
		url, code := rd.Redirect()
		http.Redirect(w, r, url, code)
		return nil
	}
	codec, _ := CodecForRequest(r, "Accept")
	data, err := codec.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", httputil.ContentType(codec.Name()))
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// DefaultErrorEncoder encodes the error to the HTTP response.
func DefaultErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
    se := errors.FromError(err)
    codec, _ := CodecForRequest(r, "Accept")
    body, err := codec.Marshal(se)
    if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
    }
    w.Header().Set("Content-Type", httputil.ContentType(codec.Name()))
    w.WriteHeader(int(se.Code))
    _, _ = w.Write(body)
}
```

需要重新写两个方法来替换默认这个两个方法, 再到项目中的 server/http.go 文件中重新注册这个两个方法就可以

### 本项目修改方式
1. 在 pkg 文件下新建 encoder/encoder.go
2. 在 encoder/encoder.go 重写 **DefaultResponseEncoder** 与 **DefaultErrorEncoder** 两个方法：
```go
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

```
3. 修改 app/service/internal/server/http.go 文件
```go

func NewHTTPServer(c *conf.Server, user *service.DemoService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.ErrorEncoder(encoder.ErrorEncoder),  // 注册返回错误的格式
		http.ResponseEncoder(encoder.ResponseEncoder),  // 注册返回正常的格式
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterDemoHTTPServer(srv, user)
	return srv
}

```
