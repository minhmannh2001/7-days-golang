package gee

import (
	"fmt"
	"strings"
)

// node đại diện cho một nút trong cây tiền tố (trie)
type node struct {
	pattern  string  // pattern là đường dẫn đầy đủ, chỉ có giá trị ở nút lá
	part     string  // part là phần của đường dẫn tương ứng với nút này
	children []*node // children chứa các nút con
	isWild   bool    // isWild = true nếu part chứa ":" hoặc "*"
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// insert thêm một route mới vào cây
// pattern: đường dẫn đầy đủ, parts: các phần của đường dẫn, height: độ cao hiện tại
func (n *node) insert(pattern string, parts []string, height int) {
	// Nếu đã duyệt hết các phần của đường dẫn, đánh dấu nút hiện tại là nút lá
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	// Tìm nút con phù hợp với phần hiện tại
	child := n.matchChild(part)
	// Nếu không tìm thấy, tạo nút con mới
	if child == nil {
		// Đánh dấu isWild = true nếu part bắt đầu bằng ":" hoặc "*"
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	// Đệ quy để thêm phần còn lại của đường dẫn
	child.insert(pattern, parts, height+1)
}

// search tìm kiếm nút phù hợp với đường dẫn
// parts: các phần của đường dẫn cần tìm, height: độ cao hiện tại
func (n *node) search(parts []string, height int) *node {
	// Điều kiện dừng: đã duyệt hết parts hoặc gặp wildcard "*"
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// Nếu nút hiện tại không phải nút lá (pattern rỗng), trả về nil
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	// Tìm tất cả các nút con có thể khớp với phần hiện tại
	children := n.matchChildren(part)

	// Duyệt qua từng nút con và tìm kiếm đệ quy
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
