package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dqn/chatlog"
	"github.com/gosuri/uilive"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var writer = uilive.New()
var printer = message.NewPrinter(language.English)

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

func printCurrencies(currencies map[string]float64) {
	keys := make([]string, 0, len(currencies))
	for k := range currencies {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(writer, "%s%s\n", key, printer.Sprintf("%.2f", currencies[key]))
	}
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
	currencies := make(map[string]float64)

	writer.Start()

	err := c.HandleChat(func(renderer chatlog.ChatRenderer) error {
		lcpmr, ok := renderer.(*chatlog.LiveChatPaidMessageRenderer)
		if !ok {
			return nil
		}

		amountStr := lcpmr.PurchaseAmountText.SimpleText

		if amountStr == "" {
			return nil
		}

		unit, amount, err := parseSuperChat(amountStr)
		if err != nil {
			return err
		}

		currencies[unit] += amount

		printCurrencies(currencies)
		time.Sleep(time.Millisecond * 25) // wait for refreshing outputs

		return nil
	})

	writer.Stop()

	return err
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}
