// Worker Pool 实现
package main

import (
	"bufio"
	"context"
	"os"
	"strings"
	"sync"
)

// 搜索任务
type SearchTask struct {
	FilePath string
}

// 搜索结果或错误
type SearchResultOrErr struct {
	Result *SearchResult
	Err    error
}

// 搜索任务的真实实现
func searchFile(ctx context.Context, task SearchTask, keyword string, ignoreCase bool) <-chan SearchResultOrErr {
	out := make(chan SearchResultOrErr)
	go func() {
		//这段代码执行结束后会执行defer这行代码
		defer close(out)
		//打开文件
		file, err := os.Open(task.FilePath)
		//判断打开文件是否有错
		if err != nil {
			out <- SearchResultOrErr{Err: err}
			return
		}
		//这段代码执行结束后会执行defer这行代码
		defer file.Close()
		//用缓冲区进行读取
		scanner := bufio.NewScanner(file)
		//按行读取
		lineNum := 0
		for scanner.Scan() {
			select {
			//要判断是否超时或者取消
			case <-ctx.Done():
				return
			default:
			}
			lineNum++
			line := scanner.Text()
			checkLine := line
			kw := keyword
			//是否有忽略大小写，全部转成小写
			if ignoreCase {
				checkLine = strings.ToLower(line)
				kw = strings.ToLower(keyword)
			}
			//判断是否包含关键字
			if strings.Contains(checkLine, kw) {
				//组转返回结果
				out <- SearchResultOrErr{
					Result: &SearchResult{
						File:    task.FilePath,
						LineNum: lineNum,
						Content: strings.TrimSpace(line),
					},
				}
			}
		}
		//如果有报错，写入报错
		if err := scanner.Err(); err != nil {
			out <- SearchResultOrErr{Err: err}
		}
	}()
	return out
}

// 启动 worker pool
func RunWorkerPool(ctx context.Context, tasks []SearchTask, keyword string, ignoreCase bool, workerCount int) <-chan SearchResultOrErr {
	resultCh := make(chan SearchResultOrErr, 100)
	var wg sync.WaitGroup

	// 创建任务队列
	taskCh := make(chan SearchTask, len(tasks))
	for _, t := range tasks {
		taskCh <- t
	}
	close(taskCh)

	// 启动 workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			//代码执行结束后会执行defer这行代码
			defer wg.Done()
			//循环任务
			for task := range taskCh {
				select {
				//有超时或者中断
				case <-ctx.Done():
					return
				default:
					//调用searchFile方法
					for res := range searchFile(ctx, task, keyword, ignoreCase) {
						select {
						//把结果写入到resultCh
						case resultCh <- res:
						//有超时或者中断
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}()
	}

	// 等待所有 worker 完成后关闭结果通道
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	return resultCh
}
