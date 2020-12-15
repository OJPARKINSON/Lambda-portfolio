package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type CoinbaseResp struct {
	Data struct {
		Currency       string `json:"currency"`
		Amount         string `json: amount`
		Native_Balance struct {
			Amount string `json: amount`
		} `json: native_balance`
		Rates struct {
			GBP string `json:GBP`
		} `json:"rates"`
	} `json:"data"`
	// Quote struct {
	// 	LatestPrice string `json:"latestPrice"`
	// } `json:"quote"`
}

func getEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func coinbaseRequest(url string) CoinbaseResp {
	timestamp := fmt.Sprintf("%v", time.Now().Unix())
	apiKey := os.Getenv("apiKey")
	apiSec := os.Getenv("apiSec")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", "https://api.coinbase.com"+url, nil)
	if err != nil {
		log.Fatal(err)
	}

	h := hmac.New(sha256.New, []byte(apiSec))
	message := timestamp + req.Method + url
	h.Write([]byte(message))
	signature := hex.EncodeToString(h.Sum(nil))

	req.Header.Add("CB-ACCESS-KEY", apiKey)
	req.Header.Add("CB-ACCESS-SIGN", signature)
	req.Header.Add("CB-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("CB-VERSION", "2015-07-22")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	respJSON := CoinbaseResp{}
	json.Unmarshal([]byte(body), &respJSON)

	return respJSON
}

func main() {
	getEnv()
	exchangeRates := coinbaseRequest("/v2/exchange-rates?currency=USD")
	xrpPrice := coinbaseRequest("/v2/prices/XRP-GBP/buy")
	portfolio := coinbaseRequest("/v2/accounts/" + os.Getenv("acountID"))

	fmt.Println("£" + portfolio.Data.Native_Balance.Amount)
	fmt.Println("£" + exchangeRates.Data.Rates.GBP)
	fmt.Println("£" + xrpPrice.Data.Amount)
}

// func iexapi() {
// 	token := os.Getenv("token")
// 	// sellPrice, ammount := 27.26, 330.0
// 	options := "&types=quote,news,chart&range=1m&last=10"
// 	resp, err := http.Get("https://cloud.iexapis.com/stable/stock/cmcsa/batch?token=" + token + options)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	respJSO := CoinbaseResp{}
// 	json.Unmarshal([]byte(body), &respJSO)

// 	fmt.Println(respJSO)
// 	// fmt.Printf("%s", body)
// }
