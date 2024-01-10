package log

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testTimes = 100000

func TestWriter(t *testing.T) {
	// 空文件名
	t.Run("empty_log_name", func(t *testing.T) {
		_, err := NewWriter("", 0, 0, 0, false)
		assert.Error(t, err)
	})

	// 默认情况，不滚动
	t.Run("default", func(t *testing.T) {
		logName := "./test.log"
		w, err := NewWriter(logName, 0, 0, 0, false)
		assert.NoError(t, err)
		log.SetOutput(w)
		for i := 0; i < testTimes; i++ {
			log.Printf("this is a test line: %d\n", i)
		}
		_ = w.Close()
		time.Sleep(time.Second)
	})

	// 按照文件大小滚动
	t.Run("size", func(t *testing.T) {
		logName := "./test.log"
		w, err := NewWriter(logName, 1, 1, 2, false)
		assert.NoError(t, err)
		log.SetOutput(w)
		for i := 0; i < testTimes; i++ {
			log.Printf("this is a test line: %d\n", i)
		}
		time.Sleep(time.Second)
		_ = w.Close()
	})

	// 压缩
	t.Run("compress", func(t *testing.T) {
		logName := "./test.log"
		w, err := NewWriter(logName, 1, 1, 2, true)
		assert.NoError(t, err)
		log.SetOutput(w)
		for i := 0; i < testTimes; i++ {
			log.Printf("this is a test line: %d\n", i)
		}
		time.Sleep(time.Second)
		_ = w.Close()
	})

	// 并行
	t.Run("concurrent", func(t *testing.T) {
		writer, err := NewWriter("./test.log", 100, 0, 0, false)
		assert.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				_, _ = writer.Write([]byte(fmt.Sprintf("this is a test line: 1-%d\n", i)))
			}
		}()

		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				_, _ = writer.Write([]byte(fmt.Sprintf("this is a test line: 2-%d\n", i)))
			}
		}()
		wg.Wait()

		_ = writer.Close()
	})
}
