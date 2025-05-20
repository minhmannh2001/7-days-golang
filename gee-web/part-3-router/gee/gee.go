package gee

import (
	"log"
	"net/http"
)

// HandlerFunc định nghĩa kiểu hàm xử lý request
type HandlerFunc func(*Context)

// Engine là đối tượng chính của framework, triển khai interface ServeHTTP
type Engine struct {
	router *router // router quản lý việc định tuyến
}

// New tạo một instance mới của Engine
func New() *Engine {
	return &Engine{router: newRouter()}
}

// addRoute là phương thức nội bộ để thêm route vào router
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	engine.router.addRoute(method, pattern, handler)
}

// GET đăng ký một handler cho HTTP GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST đăng ký một handler cho HTTP POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run khởi động HTTP server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP xử lý tất cả HTTP requests
// Triển khai interface http.Handler để Engine có thể được sử dụng với http.ListenAndServe
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Tạo context mới cho request
	c := newContext(w, req)
	// Chuyển request đến router để xử lý
	engine.router.handle(c)
}
