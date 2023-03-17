package cdek

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Size struct {
	Weight float64
	Length float64
	Width  float64
	Height float64
}

type PriceSending struct {
	TariffCode        int     `json:"tariff_code"`
	TariffName        string  `json:"tariff_name"`
	TariffDescription string  `json:"tariff_description"`
	DeliveryMode      int     `json:"delivery_mode"`
	DeliverySum       float64 `json:"delivery_sum"`
	PeriodMin         int     `json:"period_min"`
	PeriodMax         int     `json:"period_max"`
}

type CDEKClient struct {
	Account  string
	Password string
	Test     bool
	BaseURL  string
	Client   *http.Client
}

func NewCDEKClient(account string, password string, test bool) *CDEKClient {
	client := &http.Client{}
	if test {
		return &CDEKClient{
			Account:  account,
			Password: password,
			Test:     true,
			BaseURL:  "https://api.edu.cdek.ru/v2",
			Client:   client,
		}
	} else {
		return &CDEKClient{
			Account:  account,
			Password: password,
			Test:     false,
			BaseURL:  "https://api.cdek.ru/v2",
			Client:   client,
		}
	}
}

func (c *CDEKClient) Calculate(addrFrom string, addrTo string, size Size) ([]PriceSending, error) {
	reqURL := fmt.Sprintf("%s/calculator/tarifflist", c.BaseURL)

	formData := url.Values{}
	formData.Add("authLogin", c.Account)
	formData.Add("secure", c.Password)
	formData.Add("senderCityId", "44")
	formData.Add("receiverCityId", "441")
	formData.Add("tariffList", "1,10")
	formData.Add("modeId", "2")
	formData.Add("goodsWeight", fmt.Sprintf("%.3f", size.Weight))
	formData.Add("goodsLength", fmt.Sprintf("%.2f", size.Length))
	formData.Add("goodsWidth", fmt.Sprintf("%.2f", size.Width))
	formData.Add("goodsHeight", fmt.Sprintf("%.2f", size.Height))

	if c.Test {
		formData.Add("test", "1")
	}

	req, err := http.NewRequest("POST", reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Body = ioutil.NopCloser(strings.NewReader(formData.Encode()))

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Failed to calculate delivery price")
	}

	var prices []PriceSending
	err = json.NewDecoder(resp.Body).Decode(&prices)
	if err != nil {
		return nil, err
	}

	return prices, nil
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}