package bot

import (
	"errors"
	"strconv"

	gdax "github.com/preichenberger/go-gdax"
)

func GetAccount(client *gdax.Client, currency string) (*gdax.Account, error) {

	accounts, err := client.GetAccounts()

	if err != nil {
		return nil, err
	}

	for _, a := range accounts {

		if a.Currency == currency {
			return &a, nil
		}
	}

	return nil, errors.New("Strange behaviour")

}

func GetBalance(client *gdax.Client, currency string) (float64, error) {

	account, err := GetAccount(client, currency)

	if err != nil {

		accountBalance, _ := strconv.ParseFloat(account.Balance, 64)

		return accountBalance, nil

	} else {
		return 0.0, err
	}
}
