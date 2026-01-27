package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// 定义命令行参数
	port := flag.Int("port", 8000, "端口号，默认为8000")
	directory := flag.String("dir", ".", "要托管的目录，默认为当前目录")
	flag.Parse()

	// 获取绝对路径
	absDir, err := filepath.Abs(*directory)
	if err != nil {
		log.Fatalf("无法获取目录的绝对路径: %v", err)
	}

	// 检查目录是否存在
	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		log.Fatalf("目录不存在: %s", absDir)
	}

	// 创建文件服务器
	fs := http.FileServer(http.Dir(absDir))

	// 设置路由处理器
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 记录请求日志
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)

		// 使用文件服务器处理请求
		fs.ServeHTTP(w, r)
	})

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("服务器启动在 http://localhost%s\n", addr)
	fmt.Printf("托管目录: %s\n", absDir)
	fmt.Printf("按 Ctrl+C 停止服务器\n")

	// 启动服务器
	log.Fatal(http.ListenAndServe(addr, nil))
}
