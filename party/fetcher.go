package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func main() {
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
		chromedp.Nodes("#search_party", &nodes, chromedp.ByQuery),
	); err != nil {
		log.Fatal(err)
	}
	n := nodes[0]
	log.Println("nodes", n)
	err := chromedp.Run(ctx,
		chromedp.Nodes("#search_party", &nodes, chromedp.ByQuery),
		chromedp.MouseClickNode(nodes[0]),
		chromedp.Sleep(10*time.Second),
		chromedp.WaitVisible(`#123`),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("nodes", nodes)
}
