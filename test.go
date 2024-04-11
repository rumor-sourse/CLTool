package main

import (
	"context"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func main() {
	// Load HTML template from file
	data, err := ioutil.ReadFile("template.html")
	if err != nil {
		log.Fatal(err)
	}

	// Replace placeholder URL with the target URL
	html := strings.Replace(string(data), "http://example.com", "http://175.205.17.185:5531", -1)

	// Start Chrome
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()

	// Create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Navigate to the HTML template and capture a screenshot
	var buf []byte
	err = chromedp.Run(ctx, fullScreenshot(html, &buf))
	if err != nil {
		log.Fatal(err)
	}

	// Save the screenshot to a file
	if err := ioutil.WriteFile("screenshot.png", buf, 0644); err != nil {
		log.Fatal(err)
	}
}

// fullScreenshot navigates to a HTML string and captures a full screenshot.
func fullScreenshot(html string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(`data:text/html,` + html),
		chromedp.Sleep(5 * time.Second), // wait for the page to load
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Capture screenshot
			buf, err := page.CaptureScreenshot().WithQuality(90).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
