package main

import (
	"fmt"
	"log"
	"time"

	bot "bot/bot"

	gdax "github.com/preichenberger/go-gdax"
)

//https://api-public.sandbox.pro.coinbase.com

func main() {

	var order *gdax.Order

	client := setUp()

	isTrade := false

	lastPrice := 0.0

	for true {

		//Get balances
		balanceUSD := bot.GetBalance(client, "USD")
		balanceBTC := bot.GetBalance(client, "BTC")
		log.Println("USD$ " + fmt.Sprintf("%f", balanceUSD) + " --- BTC " + fmt.Sprintf("%f", balanceBTC))

		//Get BTC price
		entry := bot.GetPrice()
		log.Println("1 BTC = " + fmt.Sprintf("%f", entry.Price))

		if !isTrade {

			// Only buy if the sell order has been executed or its the first time
			if lastPrice == 0.0 || order.Settled == true {

				//Buy!
				if entry.Price > lastPrice {

					// Calculate position
					positionBTC := bot.CalculateBTCPosition(entry.Price, balanceUSD-10)

					//Place order limitted
					order = bot.BuyOrderBTC(entry.Price, positionBTC, client)

					log.Println(order.Status)

					isTrade = true
				}

			}

		} else {

			// Only sell if the buy order has been executed
			if order.Settled == true {

				//Sell!
				if entry.Price < lastPrice {

					//Place order limitted
					order = bot.SellOrderBTC(entry.Price, balanceBTC, client)

					isTrade = false

				}
			}

			// Refresh order
			order = bot.GetOrder(order.Id, client)
			log.Println("Order: " + order.Id + " --- Status: " + order.Status + " ---- Seetled? " + fmt.Sprintf("%t", order.Settled))

		}

		time.Sleep(10 * time.Second)

	}

}

func setUp() *gdax.Client {

	secret := ""
	key := ""
	passphrase := ""

	client := gdax.NewClient(secret, key, passphrase)
	client.BaseURL = "https://api-public.sandbox.pro.coinbase.com"

	return client
}
