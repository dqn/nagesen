package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/dqn/chatlog"
	"golang.org/x/text/message"
)

const format = "%s%.2f"

func splitIntoUnitAndAmount(str string) (string, float64) {
	r := []rune(str)
	for i, v := range r {
		if !unicode.IsDigit(v) {
			continue
		}
		unit := strings.TrimSpace(string(r[:i]))
		amount, err := strconv.ParseFloat(strings.ReplaceAll(string(r[i:]), ",", ""), 64)
		if err != nil {
			break
		}
		return unit, amount
	}
	return "", 0
}

func run() error {
	c, err := chatlog.New("_i_AxXSfceM")
	if err != nil {
		return err
	}

	m := make(map[string]float64)
	p := message.NewPrinter(message.MatchLanguage("en"))
	i := 0
	for c.Continuation != "" {
		resp, err := c.Fecth()
		if err != nil {
			return err
		}
		for _, continuationAction := range resp {
			for _, chatAction := range continuationAction.ReplayChatItemAction.Actions {
				amountText := chatAction.AddChatItemAction.Item.LiveChatPaidMessageRenderer.PurchaseAmountText.SimpleText
				if amountText == "" {
					continue
				}
				unit, amount := splitIntoUnitAndAmount(amountText)
				if unit == "" || amount == 0 {
					return fmt.Errorf("\ncannot parse %s\n", amountText)
				}
				m[unit] += amount
			}
		}
		fmt.Printf("\r")
		for k, v := range m {
			p.Printf(format+" ", k, v)
		}
		if i++; i > 10 {
			break
		}
	}

	fmt.Fprint(os.Stdout, "\r \r")
	for k, v := range m {
		p.Printf(format+"\n", k, v)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
