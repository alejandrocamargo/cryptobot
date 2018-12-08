package bot

import (
	"fmt"
	"log"

	gdax "github.com/preichenberger/go-gdax"
)

func CalculateBTCPosition(btcPrice float64, moneyAvailable float64) float64 {

	return moneyAvailable / btcPrice
}

/* Places  limit order */
func BuyOrderBTC(price float64, btc float64, client *gdax.Client) *gdax.Order {

	priceStr := fmt.Sprintf("%f", price)
	btcStr := fmt.Sprintf("%f", btc)

	log.Println("Placing BUY order: " + btcStr + " BTC at " + priceStr)

	order := gdax.Order{
		Price:     priceStr,
		Size:      btcStr,
		Side:      "buy",
		ProductId: "BTC-USD",
	}

	savedOrder, err := client.CreateOrder(&order)

	if err != nil {
		println(err.Error())
	}

	return &savedOrder
}

/* */
func SellOrderBTC(price float64, btc float64, client *gdax.Client) *gdax.Order {

	priceStr := fmt.Sprintf("%f", price)
	btcStr := fmt.Sprintf("%f", btc)

	log.Println("Placing SELL order: " + btcStr + " BTC at " + priceStr)

	order := gdax.Order{
		Price:     priceStr,
		Size:      btcStr,
		Side:      "sell",
		ProductId: "BTC-USD",
	}

	savedOrder, err := client.CreateOrder(&order)

	if err != nil {
		println(err.Error())
	}

	return &savedOrder
}

func GetOrder(id string, client *gdax.Client) *gdax.Order {

	order, _ := client.GetOrder(id)

	return &order

}
