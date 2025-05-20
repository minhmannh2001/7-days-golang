package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern phân tích chuỗi pattern thành các phần
// Chỉ cho phép một dấu "*" trong pattern
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			// Nếu gặp dấu "*", dừng việc phân tích (phần còn lại sẽ được bắt bởi "*")
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// getRoute tìm kiếm route phù hợp và trích xuất các tham số động
// Trả về nút phù hợp và map các tham số
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	// Tìm nút phù hợp trong cây trie
	n := root.search(searchParts, 0)

	if n != nil {
		// Nếu tìm thấy, trích xuất các tham số động từ URL
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			// Xử lý tham số động với dấu ":" (ví dụ: ":name")
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			// Xử lý tham số động với dấu "*" (ví dụ: "*filepath")
			if part[0] == '*' && len(part) > 1 {
				// Ghép tất cả các phần còn lại của URL thành một chuỗi
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
