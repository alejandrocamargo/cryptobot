/*
scp -r /Users/alejandrocamargo/go/src/bot pi@192.168.1.21:/home/pi/go/src
*/

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

	// Initialization block
	/*	order = bot.GetOrder("30fa8b1d-4850-4743-aa47-16e3710c6a38", client)
		orderID := order.Id
		isTrade = true*/
	//////////////////////////

	orderID := ""

	for true {

		//Get balances
		balanceEUR := bot.GetBalance(client, "EUR")
		balanceBTC := bot.GetBalance(client, "BTC")
		log.Println("EUR € " + fmt.Sprintf("%f", balanceEUR) + " --- BTC " + fmt.Sprintf("%f", balanceBTC))

		//Get BTC price
		entry := bot.GetPrice()
		log.Println("1 BTC = " + fmt.Sprintf("%f", entry.Price))

		if !isTrade {

			// Only buy if the sell order has been executed or its the first time
			if lastPrice == 0.0 || order.Settled == true {

				//Buy!
				if entry.Price > lastPrice {

					// Calculate position
					positionBTC := bot.CalculateBTCPosition(entry.Price, balanceEUR-10)

					//Place order limitted
					order = bot.BuyOrderBTC(entry.Price-1, positionBTC, client)

					orderID = order.Id

					isTrade = true
				}

			}

		} else {

			// Only sell if the buy order has been executed
			if order.Settled == true {

				//Sell!
				if entry.Price < lastPrice {

					// if current BTC price is bigger than order price, sell at current price
					if entry.Price > bot.ParseFloat(order.Price) {

						order = bot.SellOrderBTC(entry.Price+1, balanceBTC, client)

						//if current BTC price is lower than order price, sell at order price
					} else {

						order = bot.SellOrderBTC(bot.ParseFloat(order.Price), balanceBTC, client)

					}

					orderID = order.Id

					isTrade = false

				}
			}

		}

		// Refresh order
		order = bot.GetOrder(orderID, client)
		log.Println("Order " + order.Type + ": " + orderID + " --- Status: " + order.Status + " --- Price: " + order.Price + "€ ---- Seetled? " + fmt.Sprintf("%t", order.Settled))

		lastPrice = entry.Price

		time.Sleep(60 * time.Second)

	}

}

func setUp() *gdax.Client {

	client := gdax.NewClient(secret, key, passphrase)
	//client.BaseURL = "https://api-public.sandbox.pro.coinbase.com"
	client.BaseURL = "https://api.pro.coinbase.com"

	return client
}
