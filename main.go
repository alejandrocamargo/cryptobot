package main

import (
	"fmt"
	"log"
	"os"
	"time"

	bot "bot/bot"

	gdax "github.com/preichenberger/go-gdax"
)

//https://api-public.sandbox.pro.coinbase.com

func main() {

	var order *gdax.Order
	client := setUp()
	isTrade := false
	bearCount := 0

	lastPrice := 0.0
	orderID := ""

	if bot.ListOrders(client) {
		entry := bot.GetPrice()
		log.Println("1 BTC = " + fmt.Sprintf("%f", entry.Price))
		log.Println("Waiting 10 seconds to start!")
		time.Sleep(10 * time.Second)
	}

	// Check arguments the command itself it's one
	if len(os.Args) > 1 {
		lastPrice = bot.ParseFloat(os.Args[1])
		orderID = os.Args[2]
		isTrade = true
		order = bot.GetOrder(orderID, client)

		log.Print("Initialization:  lastPrice --> " + os.Args[1] + " orderID --> " + os.Args[2])
	} else {
		log.Print("No initialization block")
	}

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

					//Place order limited
					order = bot.BuyOrderBTC(entry.Price-0.5, positionBTC, client)

					orderID = order.Id

					isTrade = true
				}

			}

		} else {

			// Only sell if the buy order has been executed
			if order.Settled == true {

				//Sell if it goes down two times in a row
				if entry.Price < lastPrice {

					bearCount++

					log.Println("BearCount: " + fmt.Sprintf("%d", bearCount))

					if bearCount == 2 {

						// if current BTC price is bigger than order price, sell at current price
						if entry.Price > bot.ParseFloat(order.Price) {

							order = bot.SellOrderBTC(entry.Price+0.5, balanceBTC, client)

							//if current BTC price is lower than order price, sell at order price
						} else {

							order = bot.SellOrderBTC(bot.ParseFloat(order.Price), balanceBTC, client)

						}

						orderID = order.Id

						isTrade = false

						bearCount = 0
					}

				} else if entry.Price == lastPrice {
					//same price, do nothing
				} else {
					bearCount = 0
				}
			}

		}

		// Refresh order
		order = bot.GetOrder(orderID, client)
		log.Println("Order " + order.Type + ": " + orderID + " --- Status: " + order.Status + " --- Price: " + order.Price + "€ ---- Seetled? " + fmt.Sprintf("%t", order.Settled))

		lastPrice = entry.Price

		time.Sleep(10 * time.Second)

	}

}

func setUp() *gdax.Client {

	client := gdax.NewClient(secret, key, passphrase)
	//client.BaseURL = "https://api-public.sandbox.pro.coinbase.com"
	client.BaseURL = "https://api.pro.coinbase.com"

	return client
}
