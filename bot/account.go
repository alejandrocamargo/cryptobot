package bot

import (
	"strconv"

	gdax "github.com/preichenberger/go-gdax"
)

func GetAccount(client *gdax.Client, currency string) *gdax.Account {

	accounts, _ := client.GetAccounts()

	for _, a := range accounts {

		if a.Currency == currency {
			return &a
		}
	}

	return nil

}

func GetBalance(client *gdax.Client, currency string) float64 {

	account := GetAccount(client, currency)

	accountBalance, _ := strconv.ParseFloat(account.Balance, 64)

	return accountBalance
}
