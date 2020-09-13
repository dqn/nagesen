package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dqn/chatlog"
	"golang.org/x/text/message"
)

// parseSuperChat returns unit, amount and error.
func parseSuperChat(str string) (string, float64, error) {
	unit := strings.TrimRight(str, "0123456789.,")
	amountStr := strings.ReplaceAll(strings.TrimLeft(str, unit), ",", "")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return "", 0, err
	}

	return unit, amount, nil
}

func run() error {
	flag.Parse()
	flag.Usage = func() {
		fmt.Println("Usage: nagesen <video-id>")
	}

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	videoID := flag.Arg(0)
	c := chatlog.New(videoID)
	currencies := make(map[string]float64, 8) // size is inferred currency types
	p := message.NewPrinter(message.MatchLanguage("en"))

	return c.HandleChatItem(func(item *chatlog.ChatItem) error {
		amountText := item.LiveChatPaidMessageRenderer.PurchaseAmountText.SimpleText
		if amountText == "" {
			return nil
		}

		unit, amount, err := parseSuperChat(amountText)
		if err != nil {
			return err
		}

		currencies[unit] += amount

		fmt.Printf("\r")
		for k, v := range currencies {
			var format string
			if k == "¥" {
				format = "%s%.0f"
			} else {
				format = "%s%.2f"
			}
			p.Printf(format+" ", k, v)
		}

		return nil
	})
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}
