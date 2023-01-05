package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/extrame/xls"
)

type Party struct {
	id                  int
	name                string
	chairman            string
	established_date    time.Time
	filing_date         time.Time
	main_office_address string
	mailing_address     string
	phone_number        string
	status              int
}

func main() {
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

	done := make(chan string, 1)
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if ev, ok := ev.(*browser.EventDownloadProgress); ok {
			if ev.State == browser.DownloadProgressStateCompleted {
				done <- ev.GUID
				close(done)
			}
		}

	})

	var ns []*cdp.Node
	if err := chromedp.Run(ctx,
		chromedp.MouseClickNode(nodes[0]),
		chromedp.WaitVisible("ContentPlaceHolder1_BTN_Export_Excel", chromedp.ByID),
		chromedp.Nodes("ContentPlaceHolder1_BTN_Export_Excel", &ns, chromedp.ByID),
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).WithDownloadPath(tmpDir).WithEventsEnabled(true),
	); err != nil {
		log.Println(err)
	}

	filename := ""
	if err := chromedp.Run(ctx,
		chromedp.MouseClickNode(ns[0]),
		chromedp.ActionFunc(func(ctx context.Context) error {
			filename = <-done
			return nil
		}),
	); err != nil {
		log.Println(err)
	}

	parties := []Party{}

	if xlFile, err := xls.Open(tmpDir+"/"+filename, "utf-8"); err == nil {
		if sheet := xlFile.GetSheet(0); sheet != nil {
			for row := 0; row <= int(sheet.MaxRow); row++ {
				parties = append(parties, rowToParty(sheet.Row(row)))
			}
		}
	}

	fmt.Println(parties)
}

func rowToParty(row *xls.Row) Party {
	id, _ := strconv.Atoi(row.Col(0))
	return Party{
		id:                  id,
		name:                row.Col(1),
		chairman:            row.Col(2),
// 		established_date:    row.Col(3),
// 		filing_date:         row.Col(4),
		main_office_address: row.Col(5),
		mailing_address:     row.Col(6),
		phone_number:        row.Col(7),
		status:              statusStrToNum(row.Col(8)),
	}
}

func statusStrToNum(status string) int {
	switch status {
	case "一般":
		return 1
	case "撤銷備案":
		return 2
	case "自行解散":
		return 3
	case "失聯":
		return 4
	case "廢止備案":
		return 5
	}

	return 0
}
