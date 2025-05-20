package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

// Context chứa thông tin về request và response
type Context struct {
	// Đối tượng gốc
	Writer http.ResponseWriter
	Req    *http.Request
	// Thông tin request
	Path   string
	Method string
	Params map[string]string // Lưu trữ các tham số động từ URL
	// Thông tin response
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// Param trả về giá trị của tham số động trong URL
// Ví dụ: với route "/user/:id", Param("id") sẽ trả về giá trị thực tế
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON thiết lập response dạng JSON
// Mã hóa đối tượng obj thành JSON và gửi về client
func (c *Context) JSON(code int, obj interface{}) {
	// Thiết lập header Content-Type là application/json
	c.SetHeader("Content-Type", "application/json")
	// Thiết lập mã trạng thái HTTP
	c.Status(code)
	// Sử dụng json.NewEncoder để mã hóa obj thành JSON
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		// Xử lý lỗi nếu không thể mã hóa
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
