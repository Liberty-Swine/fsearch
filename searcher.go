package main

import (
	"os"
	"path/filepath"
	"strings"
)

// 收集指定目录下所有符合扩展名过滤的文件路径
func collectFiles(rootDir string, extFilter string) ([]SearchTask, error) {
	var tasks []SearchTask
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 忽略无法访问的文件
		}
		if info.IsDir() {
			return nil
		}
		// 扩展名过滤
		if extFilter != "" && !strings.HasSuffix(path, extFilter) {
			return nil
		}
		tasks = append(tasks, SearchTask{FilePath: path})
		return nil
	})
	return tasks, err
}
