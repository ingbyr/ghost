# 静态文件托管服务器

这是一个简单的静态文件托管服务器，功能类似于 `python -m http.server`。

## 功能特性

- 托管指定目录的文件
- 显示目录列表
- 提供文件下载服务
- 记录访问日志

## 使用方法

### 运行服务器

```bash
# 使用默认端口8000，托管当前目录
go run file_server.go

# 指定端口和目录
go run file_server.go -port 8080 -dir ./mock

# 托管特定目录
go run file_server.go -dir /path/to/directory
```

### 参数说明

- `-port`: 指定端口号，默认为8000
- `-dir`: 指定要托管的目录，默认为当前目录(.)

## 示例

```bash
# 托管mock目录，使用端口8080
go run file_server.go -port 8080 -dir .

# 然后在浏览器中访问 http://localhost:8080
```

## 注意事项

- 确保指定的目录存在
- 确保指定的端口未被占用
- 此服务器仅用于开发和测试用途，不适用于生产环境