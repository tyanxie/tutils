package tutils

import (
	"encoding/json"
	"io"
)

// Close 静默批量关闭
func Close(closers ...io.Closer) {
	for _, closer := range closers {
		if closer != nil {
			_ = closer.Close()
		}
	}
}

// Json 将v使用json.Marshal转换，忽略错误
func Json(v interface{}) string {
	j, _ := json.Marshal(v)
	return string(j)
}
