package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func main() {
	//create a tmp folder for download
	tmpDir, err := ioutil.TempDir("", "chromedp-")
	if err != nil {
		log.Fatal(err)
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	var nodes []*cdp.Node
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://party.moi.gov.tw"),
		chromedp.Nodes("search_party", &nodes, chromedp.ByID),
	); err != nil {
		log.Println(err)
	}

	done := make(chan struct{})
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if ev, ok := ev.(*browser.EventDownloadProgress); ok {
			if ev.State == browser.DownloadProgressStateCompleted {
				done <- struct{}{}
				close(done)
			}
		}
	})

	var ns []*cdp.Node
	if err := chromedp.Run(ctx,
		chromedp.MouseClickNode(nodes[0]),
		chromedp.WaitVisible(`ContentPlaceHolder1_BTN_Export_Ods`, chromedp.ByID),
		chromedp.Nodes("ContentPlaceHolder1_BTN_Export_Ods", &ns, chromedp.ByID),
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).WithDownloadPath(tmpDir+"/fff.ods").WithEventsEnabled(true),
	); err != nil {
		log.Println(err)
	}

	if err := chromedp.Run(ctx,
		chromedp.MouseClickNode(ns[0]),
		chromedp.ActionFunc(func(ctx context.Context) error {
			<-done
			return nil
		}),
	); err != nil {
		log.Println(err)
	}
}
