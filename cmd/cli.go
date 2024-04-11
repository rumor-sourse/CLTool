// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.
package main

import (
	"context"
	"flag"
	"github.com/chromedp/chromedp"
	"log"
	"os"
)

func main() {
	// 获取url
	var url string
	flag.StringVar(&url, "url", "", "URL to capture screenshot")
	flag.Parse()

	if url == "" {
		log.Fatal("URL must be provided")
	}

	// 创建context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	// capture screenshot
	var buf []byte

	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(url, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("fullScreenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	log.Printf("wrote fullScreenshot.png")
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Note: chromedp.FullScreenshot overrides the device's emulation settings. Use
// device.Reset to reset the emulation and viewport settings.
func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
	}
}
