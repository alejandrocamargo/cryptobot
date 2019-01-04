package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetPrice() (Data, error) {

	response, err := http.Get("https://api.pro.coinbase.com/products/BTC-EUR/ticker")

	if err != nil {

		var data Data
		return data, err

	}

	//response, _ := http.Get("https://api-public.sandbox.pro.coinbase.com/products/BTC-USD/ticker")

	responseData, _ := ioutil.ReadAll(response.Body)

	return unMarshall(responseData), nil

}

func unMarshall(data []byte) Data {

	var entry Data

	json.Unmarshal(data, &entry)

	log.Println("1 BTC = " + fmt.Sprintf("%f", entry.Price))

	return entry

}

type Data struct {
	TradeID string  `json:",trade_id"`
	Price   float64 `json:",string"`
	Size    string  `json:",size"`
	Time    string  `json:",time"`
}
