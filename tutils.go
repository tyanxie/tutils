package tutils

import (
	"encoding/json"
	"io"
	"os"
)

// Close 静默批量关闭
func Close(closers ...io.Closer) {
	for _, closer := range closers {
		if closer != nil {
			_ = closer.Close()
		}
	}
}

// Remove 静默批量删除文件
func Remove(names ...string) {
	for _, name := range names {
		if name != "" {
			_ = os.Remove(name)
		}
	}
}

// Json 将v使用json.Marshal转换，忽略错误
func Json(v interface{}) string {
	j, _ := json.Marshal(v)
	return string(j)
}

// Ternary 三元表达式，当condition成立时返回x，否则返回y
func Ternary[T interface{}](condition bool, x, y T) T {
	if condition {
		return x
	}
	return y
}
