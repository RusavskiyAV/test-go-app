package currencyConvertor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	account "test3/src/model/account/entity"
)

type Interface interface {
	Convert(money *account.Money, currency string) (*account.Money, error)
}

type CurrencyConverter struct {
	APIKey string
}

type Response struct {
	Result  float64 `json:"result"`
	Success bool    `json:"success"`
}

func (cc *CurrencyConverter) Convert(money *account.Money, currency string) (*account.Money, error) {
	const url = "https://api.apilayer.com/currency_data/convert"

	requestUrl := fmt.Sprintf("%s?to=%s&from=%s&amount=%f", url, currency, money.Currency().Code, money.AsMajorUnits())

	client := &http.Client{}
	req, err := http.NewRequest("GET", requestUrl, nil)
	req.Header.Set("apikey", cc.APIKey)

	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Do(req)

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)

	var resp Response

	if err = json.Unmarshal(body, &resp); err != nil {
		return &account.Money{}, err
	}

	newMoney, err := account.NewMoneyFromFloatWithCurrency(resp.Result, currency)

	if err != nil {
		return &account.Money{}, err
	}

	return newMoney, nil
}
