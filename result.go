// 搜索结果结构

package main

type SearchResult struct {
	File    string
	LineNum int
	Content string
}

// 用于按文件路径排序
type ByFile []SearchResult

func (s ByFile) Len() int {
	return len(s)
}

func (s ByFile) Less(i, j int) bool {
	return s[i].File < s[j].File
}

func (s ByFile) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
