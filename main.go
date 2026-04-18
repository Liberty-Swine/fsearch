package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"sort"
	"time"
)

func main() {
	// 命令行参数
	var (
		dir        string
		keyword    string
		workers    int
		ignoreCase bool
		ext        string
		timeout    time.Duration
	)
	//定义参数
	flag.StringVar(&dir, "dir", ".", "搜索目录")
	flag.StringVar(&keyword, "keyword", "", "搜索关键字（必填）")
	flag.IntVar(&workers, "workers", runtime.NumCPU(), "并发 worker 数量")
	flag.BoolVar(&ignoreCase, "ignore-case", false, "忽略大小写")
	flag.StringVar(&ext, "ext", "", "文件扩展名过滤（如 .go）")
	flag.DurationVar(&timeout, "timeout", 0, "超时时间（如 30s）")
	flag.Parse()

	if keyword == "" {
		fmt.Fprintln(os.Stderr, "错误：必须指定 -keyword")
		flag.Usage()
		os.Exit(1)
	}

	// 1. 收集文件列表
	fmt.Println("正在收集文件...")
	tasks, err := collectFiles(dir, ext)
	if err != nil {
		fmt.Fprintf(os.Stderr, "收集文件失败: %v\n", err)
		os.Exit(1)
	}
	if len(tasks) == 0 {
		fmt.Println("没有找到符合条件的文件")
		return
	}

	fmt.Printf("找到 %d 个文件，使用 %d 个 worker 并发搜索...\n", len(tasks), workers)

	// 2. 创建带超时的 context
	ctx := context.Background()
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	// 3. 启动 worker pool 搜索
	resultCh := RunWorkerPool(ctx, tasks, keyword, ignoreCase, workers)

	// 4. 收集结果
	var results []SearchResult
	var errors []error
	for res := range resultCh {
		if res.Err != nil {
			errors = append(errors, res.Err)
		} else if res.Result != nil {
			results = append(results, *res.Result)
		}
	}

	// 5. 检查是否超时
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("\n搜索超时，已返回部分结果")
	}

	// 6. 输出结果
	if len(results) == 0 {
		fmt.Println("未找到匹配内容")
	} else {
		sort.Sort(ByFile(results))
		fmt.Printf("\n找到 %d 处匹配:\n", len(results))
		for _, r := range results {
			fmt.Printf("%s (第%d行): %s\n", r.File, r.LineNum, r.Content)
		}
	}

	// 7. 输出错误汇总（可选）
	if len(errors) > 0 {
		fmt.Printf("\n发生 %d 个错误（如文件无法读取）\n", len(errors))
	}
}
