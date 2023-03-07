package main

import (
	"context"
	"fmt"
	"github.com/aiviaio/go-binance/v2"
	"strings"
)

func main() {
	var (
		apiKey    = "your api key"
		secretKey = "your secret key"
		ctx       = context.Background()
	)
	ch := make(chan map[string]string)

	client := binance.NewClient(apiKey, secretKey)

	resp, err := client.NewExchangeInfoService().Do(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	var coins []string
	for _, s := range resp.Symbols {
		if strings.Contains(s.Symbol, "USDT") {
			coins = append(coins, s.Symbol)
		}
		if len(coins) == 5 {
			break
		}
	}

	for _, c := range coins {
		go func(symbol string) {
			result := make(map[string]string)
			prices, err := client.NewListPricesService().Symbol(symbol).Do(ctx)
			if err != nil {
				fmt.Println(err)
				return
			}

			result[symbol] = prices[0].Price

			ch <- result
		}(c)
	}

	for {
		select {
		case pair, ok := <-ch:
			if !ok {
				return
			}
			for s, v := range pair {
				fmt.Println(s + " " + v)
			}
		}
	}
}
