package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	// 获取url 和 proxy
	var url, proxy string
	flag.StringVar(&url, "url", "", "URL to capture screenshot")
	flag.StringVar(&proxy, "proxy", "", "Proxy server address")
	flag.Parse()

	if url == "" {
		log.Fatal("URL must be provided")
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.WindowSize(1920, 1080),
	)

	if proxy != "" {
		opts = append(opts, chromedp.ProxyServer(proxy))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	// 创建context
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	//	图片存储的位置
	var buf []byte

	// 运行chromedp组件，截图
	if err := chromedp.Run(ctx, fullScreenshot(url, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("fullScreenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	log.Printf("wrote fullScreenshot.png")

	// 创建一个新的文件
	file, err := os.Create("template.html")
	if err != nil {
		log.Fatal(err)
	}

	html := `<!DOCTYPE html>
<html>
<head>
    <style>
        #url-bar {
            height: 30px;
            background: #eee;
            padding: 5px;
            font-family: monospace;
        }
        #content {
            position: absolute;
            top: 40px;
            bottom: 0;
            left: 0;
            right: 0;
        }
        img {
            max-width: 100%;
            height: auto;
        }
    </style>
</head>
<body>
<div id="url-bar"></div>
<div id="content">
    <img src="./fullScreenshot.png" alt="Image"> <!-- replace with the actual image URL -->
</div>

<script>
    var url = 'http://example.com'; // replace with the target URL
    document.getElementById('url-bar').textContent = url;
</script>
</body>
</html>
`
	html = strings.Replace(html, "http://example.com", url, -1)
	// 写入HTML内容
	_, err = file.WriteString(html)
	if err != nil {
		log.Fatal(err)
	}

	// 确保所有的内容都被写入文件
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}

	err = file.Close()
	if err != nil {
		return
	}
	//第二次截图
	go screen(ctx, "http://localhost:3000/template.html")
	fs := http.FileServer(http.Dir("."))
	srv := &http.Server{
		Addr:    ":3000",
		Handler: fs,
	}
	// 启动服务器
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// 如果监听失败，则打印错误并退出
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// 在10秒后关闭服务器
	time.AfterFunc(10*time.Second, func() {
		// 创建一个有超时的上下文，以便关闭服务器
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			// 如果关闭服务器失败，则打印错误
			log.Printf("Server Shutdown: %v", err)
		}

		log.Println("Server exited")
		os.Exit(0)
	})

	// 等待服务器关闭
	select {}
}

// fullScreenshot 截图
func fullScreenshot(url string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.FullScreenshot(res, quality),
	}
}
func screen(ctx context.Context, url string) {
	time.Sleep(5 * time.Second)
	var buf []byte
	// 截图
	if err := chromedp.Run(ctx, fullScreenshot(url, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	//输出图片
	if err := os.WriteFile("Screenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	log.Printf("wrote Screenshot.png")
	//删除中间文件
	err := os.Remove("fullScreenshot.png")
	if err != nil {
		fmt.Println(err)
	}
	err = os.Remove("template.html")
	if err != nil {
		fmt.Println(err)
	}
}
