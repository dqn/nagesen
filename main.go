package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/dqn/chatlog"
)

func parseYen(str string) int {
	r := []rune(str)
	if r[0] != '￥' {
		fmt.Printf("warning: cannot parse %s\n", str)
		return 0
	}
	i, err := strconv.Atoi(strings.Replace(string(r[1:]), ",", "", -1))
	if err != nil {
		panic(err)
	}
	return i
}

func run() error {
	c, err := chatlog.New("_i_AxXSfceM")
	if err != nil {
		return err
	}
	amount := 0
	for c.Continuation != "" {
		resp, err := c.Fecth()
		if err != nil {
			return err
		}
		for _, continuationAction := range resp {
			for _, chatAction := range continuationAction.ReplayChatItemAction.Actions {
				if amountText := chatAction.AddChatItemAction.Item.LiveChatPaidMessageRenderer.PurchaseAmountText.SimpleText; amountText != "" {
					amount += parseYen(amountText)
					fmt.Println(amountText)
				}
			}
		}
	}
	fmt.Printf("\ntotal: %d\n", amount)
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
