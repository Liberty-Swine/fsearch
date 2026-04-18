# fsearch – 并发文件内容搜索工具

一个用 Go 编写的命令行工具，递归搜索目录下所有文件中包含关键字的行，支持并发、超时、扩展名过滤等功能。

## 功能特性

- 并发搜索（worker pool），充分利用多核 CPU
- 支持超时控制（`-timeout`）
- 支持忽略大小写（`-ignore-case`）
- 支持按文件扩展名过滤（`-ext`）
- 搜索结果按文件路径排序输出

## 安装

```bash
git clone https://github.com/Liberty-Swine/fsearch.git
cd fsearch
go build -o fsearch.exe
```

## 使用方法

```bash
./fsearch -dir ./src -keyword "func" -workers 8 -ignore-case -ext .go -timeout 30s
```

### 参数说明

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-dir` | 搜索目录 | `.` |
| `-keyword` | 搜索关键字（必填） | 无 |
| `-workers` | 并发 worker 数量 | CPU 核心数 |
| `-ignore-case` | 忽略大小写 | false |
| `-ext` | 文件扩展名过滤（如 `.go`） | 空（不过滤） |
| `-timeout` | 超时时间（如 `30s`） | 0（无超时） |

## 示例

```bash
# 在当前目录下搜索包含 "error" 的所有文件
./fsearch -keyword "error"

# 在 ./logs 目录下搜索 "panic"，忽略大小写，只搜索 .log 文件，超时 10 秒
./fsearch -dir ./logs -keyword "panic" -ignore-case -ext .log -timeout 10s
```

## 项目结构

```
fsearch/
├── main.go       # 入口，参数解析
├── searcher.go   # 目录遍历
├── worker.go     # worker pool 和搜索逻辑
├── result.go     # 结果结构及排序
└── go.mod
```