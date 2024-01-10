package log

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"go.uber.org/atomic"
)

const (
	backupTimeFormat = "bk-20060102-150405.00000" // 备份文件后缀
	compressSuffix   = ".gz"                      // 压缩文件后缀
)

// fileInfo 文件基本信息，包含了文件创建时间
type fileInfo struct {
	timestamp time.Time
	os.FileInfo
}

// Writer 日志文件写入器
type Writer struct {
	dir        string // 文件所在目录
	filename   string // 文件路径
	maxSize    int64  // 文件最大大小
	maxAge     int64  // 文件最大保存时间
	maxBackups int64  // 文件最大数量
	compress   bool   // 是否压缩文件

	size     *atomic.Int64            // 文件当前大小
	file     *atomic.Pointer[os.File] // 文件实例
	openTime *atomic.Int64            // 文件打开时间

	mu                   sync.Mutex    // 读写锁
	notifyCleanFileCh    chan struct{} // 文件清除器等待管道
	notifyCleanFilesOnce sync.Once     // 文件清除器等待管道初始化
}

// NewWriter 创建一个写入器
func NewWriter(filename string, maxSize, maxAge, maxBackups int64, compress bool) (*Writer, error) {
	// 文件名校验
	if filename == "" {
		return nil, errors.New("invalid file path")
	}
	// 获取文件目录并创建
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}
	// 创建文件并返回
	return &Writer{
		dir:        dir,
		filename:   filename,
		maxSize:    maxSize * 1024 * 1024,
		maxAge:     maxAge,
		maxBackups: maxBackups,
		compress:   compress,
		size:       atomic.NewInt64(0),
		file:       atomic.NewPointer[os.File](nil),
		openTime:   atomic.NewInt64(0),
	}, nil
}

// Write 向文件中写入内容
func (w *Writer) Write(v []byte) (n int, err error) {
	// 每10s尝试重新打开文件
	if w.file.Load() == nil || time.Now().Unix()-w.openTime.Load() > 10 {
		w.mu.Lock()
		// 加锁后二次判断
		if w.file.Load() == nil || time.Now().Unix()-w.openTime.Load() > 10 {
			w.reopen(true)
		}
		w.mu.Unlock()
	}
	// 获取当前文件
	file := w.file.Load()
	// 如果文件依然为nil，报错返回
	if file == nil {
		return 0, errors.New("open file fail")
	}
	// 写入文件内容并记录大小
	n, err = file.Write(v)
	w.size.Add(int64(n))
	// 如果超过文件最大大小，备份文件并重新打开新文件
	if w.maxSize > 0 && w.size.Load() >= w.maxSize {
		w.mu.Lock()
		// 获得锁后二次判断
		if w.maxSize > 0 && w.size.Load() >= w.maxSize {
			w.backupFile()
		}
		w.mu.Unlock()
	}
	return n, err
}

// Close 关闭当前文件
func (w *Writer) Close() error {
	// 获取当前文件
	file := w.file.Load()
	if file == nil {
		return nil
	}
	// 关闭当前文件
	err := file.Close()
	// 设置当前文件为nil
	w.file.Store(nil)

	if w.notifyCleanFileCh != nil {
		close(w.notifyCleanFileCh)
		w.notifyCleanFileCh = nil
	}

	return err
}

// doReopenFile 重新打开文件，交给上层加锁和判断
// @param path 打开的文件路径
// @param needClose 是否需要关闭当前文件，如果为false，则上层函数必须保证w.currFile被关闭
func (w *Writer) reopen(needClose bool) {
	// 记录打开时间
	w.openTime.Store(time.Now().Unix())
	// 外层判断是否要关闭，可以少一次atomic.Load
	if needClose {
		if file := w.file.Load(); file != nil {
			_ = file.Close()
		}
	}
	// 打开新文件
	file, err := os.OpenFile(w.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	// 设置新文件
	w.file.Store(file)
	// 设置文件大小
	if st, _ := os.Stat(w.filename); st != nil {
		w.size.Store(st.Size())
	}
}

// backupFile 备份当前文件并打开新的文件，交给上层加锁和判断文件大小
func (w *Writer) backupFile() {
	// 重新设置当前文件大小
	w.size.Store(0)
	// 先关闭文件，否则windows系统下重命名文件会报错文件被占用
	if file := w.file.Load(); file != nil {
		_ = file.Close()
	}
	// 重命名文件
	backupFilename := w.filename + "." + time.Now().Format(backupTimeFormat)
	if _, e := os.Stat(w.filename); !os.IsNotExist(e) {
		_ = os.Rename(w.filename, backupFilename)
	}
	// 重新打开新的文件，因为之前已经关闭了当前文件，因此无需再reopen中再次关闭文件
	w.reopen(false)
	// 通知文件清除器
	w.notifyCleanFiles()
}

// notifyCleanFiles 通知文件清除器启动
func (w *Writer) notifyCleanFiles() {
	w.notifyCleanFilesOnce.Do(func() {
		w.notifyCleanFileCh = make(chan struct{}, 1)
		go w.runCleanFiles()
	})
	select {
	case w.notifyCleanFileCh <- struct{}{}:
	default:
	}
}

// runCleanFiles 执行文件清除
func (w *Writer) runCleanFiles() {
	for range w.notifyCleanFileCh {
		if w.maxBackups <= 0 && w.maxAge <= 0 && !w.compress {
			continue
		}
		w.cleanFiles()
	}
}

// cleanFiles 清除和压缩文件
func (w *Writer) cleanFiles() {
	// 获取文件列表
	files, err := w.getOldFiles()
	if err != nil || len(files) == 0 {
		return
	}
	// 过滤要删除的文件
	var remove, compress []fileInfo
	remove, files = w.filterNeedRemove(files)
	// 过滤需要压缩的文件
	compress = w.filterNeedCompress(files)
	// 删除文件
	w.removeFiles(remove)
	// 压缩文件
	w.compressFiles(compress)
}

// getOldFiles 按照修改时间返回旧文件 ordered by modified time.
func (w *Writer) getOldFiles() ([]fileInfo, error) {
	// 获取目录下的所有文件
	files, err := os.ReadDir(w.dir)
	if err != nil {
		return nil, fmt.Errorf("read log file directory failed: %w", err)
	}
	// 循环匹配文件
	matchFileInfos := make([]fileInfo, 0)
	filename := filepath.Base(w.filename)
	for _, f := range files {
		// 跳过目录
		if f.IsDir() {
			continue
		}
		// 获取文件信息
		info, err := f.Info()
		if err != nil {
			continue
		}
		// 获取文件修改时间
		modTime, ok := w.matchFile(f.Name(), filename)
		if !ok {
			continue
		}
		// 保存文件信息
		matchFileInfos = append(matchFileInfos, fileInfo{modTime, info})
	}
	// 按照文件时间排序
	sort.Slice(matchFileInfos, func(i, j int) bool {
		return matchFileInfos[i].timestamp.After(matchFileInfos[j].timestamp)
	})
	return matchFileInfos, nil
}

// matchFile 尝试匹配文件名称，如果匹配成功，则返回文件创建时间
func (w *Writer) matchFile(filename, prefix string) (time.Time, bool) {
	// 排除文件本身
	// a.log
	// a.log.20200712
	if filepath.Base(w.filename) == filename {
		return time.Time{}, false
	}

	// 排除非filename文件前缀
	// a.log -> a.log.20200712-1232/a.log.20200712-1232.gz
	// a.log.20200712 -> a.log.20200712.20200712-1232/a.log.20200712.20200712-1232.gz
	if !strings.HasPrefix(filename, prefix) {
		return time.Time{}, false
	}

	// 返回文件修改时间
	if st, _ := os.Stat(filepath.Join(w.dir, filename)); st != nil {
		return st.ModTime(), true
	}
	return time.Time{}, false
}

// filterNeedRemove 过滤需要删除的文件
func (w *Writer) filterNeedRemove(files []fileInfo) (remove, remaining []fileInfo) {
	// 按照文件保存最大数量过滤
	if w.maxBackups > 0 && len(files) >= int(w.maxBackups) {
		// 保留文件表
		preserved := make(map[string]struct{})
		for _, f := range files {
			// 移除文件后缀进行保留
			filename := strings.TrimSuffix(f.Name(), compressSuffix)
			preserved[filename] = struct{}{}
			// 如果当前保留的文件数量已经大与最大文件数量，直接删除
			if len(preserved) > int(w.maxBackups) {
				remove = append(remove, f)
			} else {
				remaining = append(remaining, f)
			}
		}
	}
	// 按照文件最大保存时间过滤文件
	if w.maxAge > 0 {
		// 重制文件列表
		files = remaining
		remaining = nil
		// 获取文件列表
		diff := time.Duration(int64(24*time.Hour) * w.maxAge)
		cutoff := time.Now().Add(-1 * diff)
		// 这里需要遍历剩余文件列表
		for _, f := range files {
			if f.timestamp.Before(cutoff) {
				remove = append(remove, f)
			} else {
				remaining = append(remaining, f)
			}
		}
	}
	return
}

// filterNeedCompress 过滤所有需要压缩的文件
func (w *Writer) filterNeedCompress(files []fileInfo) (compress []fileInfo) {
	if !w.compress {
		return nil
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), compressSuffix) {
			compress = append(compress, f)
		}
	}
	return
}

// removeFiles 删除文件
func (w *Writer) removeFiles(remove []fileInfo) {
	for _, f := range remove {
		_ = os.Remove(filepath.Join(w.dir, f.Name()))
	}
}

// compressFiles 压缩文件
func (w *Writer) compressFiles(compress []fileInfo) {
	for _, f := range compress {
		filename := filepath.Join(w.dir, f.Name())
		compressFilename := filename + compressSuffix
		_ = compressFile(filename, compressFilename)
	}
}

// compressFile 压缩文件，压缩完成后删除原文件
func compressFile(filename, compressFilename string) (err error) {
	// 打开原始文件
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open file failed: %w", err)
	}
	defer func() {
		// 关闭文件
		_ = f.Close()
		// 如果没有发生错误，则删除文件
		if err == nil {
			_ = os.Remove(filename)
		}
	}()
	// 创建压缩文件
	gzf, err := os.Create(compressFilename)
	if err != nil {
		return fmt.Errorf("open compressed file failed: %w", err)
	}
	defer func() {
		// 关闭压缩文件
		_ = gzf.Close()
		// 如果发生错误，则删除压缩文件
		if err != nil {
			_ = os.Remove(compressFilename)
		}
	}()
	// 打开压缩器
	gz := gzip.NewWriter(gzf)
	defer func() {
		_ = gz.Close()
	}()
	// 压缩文件
	_, err = io.Copy(gz, f)
	if err != nil {
		return fmt.Errorf("compress file failed: %w", err)
	}
	return nil
}
