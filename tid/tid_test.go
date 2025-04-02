package tid

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	fmt.Println(Generate())

	begin := time.Now()
	for range 0xffff {
		_ = Generate()
	}
	fmt.Println("time cost: ", time.Now().Sub(begin))
}

func BenchmarkGenerate(b *testing.B) {
	for b.Loop() {
		_ = Generate()
	}
}
