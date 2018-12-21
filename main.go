package main

import (
	"errors"
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
	bearCount := 0
	lastPrice := bot.GetPrice().Price
	orderID := ""

	// Look for open order an show warning
	if bot.ListOrders(client) {
		log.Println("1 BTC = " + fmt.Sprintf("%f", lastPrice))
		log.Println("Waiting 10 seconds to start!")
		time.Sleep(10 * time.Second)
	}

	// Check arguments the command itself it's one
	if len(os.Args) > 1 {
		lastPrice = bot.ParseFloat(os.Args[1])
		orderID = os.Args[2]

		order, _ = bot.GetOrder(orderID, client)

		log.Print("Initialization:  lastPrice --> " + os.Args[1] + " orderID --> " + os.Args[2])
	} else {
		log.Print("No initialization block")
	}

	for true {

		//Get balances
		balanceEUR, balanceBTC, err := getBalances(client)

		if err != nil {
			log.Fatal(err)
			time.Sleep(10 * time.Second)
			continue
		}

		//Get BTC price
		entry := bot.GetPrice()

		if order == nil || order.Side == "sell" {

			// Only buy if the sell order has been executed or there is no order at all
			if order == nil || order.Settled == true {

				//Buy!
				if entry.Price > lastPrice {

					// Calculate position
					positionBTC := bot.CalculateBTCPosition(entry.Price, balanceEUR-10)

					//Place order limited, 5€ cheaper
					order = bot.BuyOrderBTC(entry.Price-5, positionBTC, client)

					orderID = order.Id
				}

			}

		} else if order.Side == "buy" {

			// Only sell if the buy order has been executed
			if order.Settled == true {

				//Sell if it goes down two times in a row
				if entry.Price < lastPrice {

					bearCount++

					log.Println("BearCount: " + fmt.Sprintf("%d", bearCount))

					if bearCount == 2 {

						// if current BTC price is bigger than order price, sell at current price
						if entry.Price > bot.ParseFloat(order.Price) {

							order = bot.SellOrderBTC(entry.Price+1, balanceBTC, client)

							//if current BTC price is lower than order price, sell at order price
						} else {

							order = bot.SellOrderBTC(bot.ParseFloat(order.Price), balanceBTC, client)

						}

						orderID = order.Id

						bearCount = 0
					}

				} else if entry.Price == lastPrice {
					//same price, do nothing
				} else {
					bearCount = 0
				}
			}

		}

		order = refreshOrder(order, orderID, client)

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

func getBalances(client *gdax.Client) (balanceEUR float64, balanceBTC float64, err error) {

	balanceEUR, err1 := bot.GetBalance(client, "EUR")
	balanceBTC, err2 := bot.GetBalance(client, "BTC")

	if err1 != nil || err2 != nil {

		return 0.0, 0.0, errors.New("Cannot retrieve balances!")

	} else {

		log.Println("EUR € " + fmt.Sprintf("%f", balanceEUR) + " --- BTC " + fmt.Sprintf("%f", balanceBTC))

		return balanceEUR, balanceBTC, nil

	}

}

func refreshOrder(order *gdax.Order, orderID string, client *gdax.Client) *gdax.Order {

	if order != nil {

		// Refresh order
		orderP, err := bot.GetOrder(orderID, client)

		// re-assign only if no problem
		if err == nil {
			order = orderP
		}

		log.Println("Order " + order.Type + ": " + orderID + " --- Status: " + order.Status + " --- Price: " + order.Price + "€ ---- Seetled? " + fmt.Sprintf("%t", order.Settled))

		return order

	} else {

		log.Println("No order placed.")

		return nil
	}

}
