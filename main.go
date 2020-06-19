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

func parseSuperChat(str string) (string, float64, error) {
	unit := strings.TrimRight(str, "0123456789.,")
	s := strings.TrimLeft(str, unit)
	s = strings.ReplaceAll(s, ",", "")
	amount, err := strconv.ParseFloat(s, 64)
	unit = strings.TrimSpace(strings.ReplaceAll(unit, "￥", "¥"))

	return unit, amount, err
}

func run() error {
	flag.Parse()
	flag.Usage = func() {
		fmt.Println("Usage: nagesen <video-id>")
	}
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	c := chatlog.New(flag.Arg(0))
	m := make(map[string]float64, 16)
	p := message.NewPrinter(message.MatchLanguage("en"))

	err := c.HandleChatItem(func(item *chatlog.ChatItem) error {
		amountText := item.LiveChatPaidMessageRenderer.PurchaseAmountText.SimpleText
		if amountText == "" {
			return nil
		}

		unit, amount, err := parseSuperChat(amountText)
		if err != nil {
			return err
		}

		m[unit] += amount

		fmt.Printf("\r")
		for k, v := range m {
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
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}
