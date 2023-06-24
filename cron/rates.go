package cron

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"moma-api/db"
	"net/http"
	"time"
)

func RefreshRates(ticker *time.Ticker, done chan bool) {
	go func() {
		err := refreshRates()
		if err != nil {
			fmt.Printf("refresh rates failed due to: %v \n", err)
		}

		for {
			select {
			case <-done:
				ticker.Stop()
				fmt.Println("Ticker stopped")
				return
			case <-ticker.C:
				err := refreshRates()
				if err != nil {
					fmt.Printf("refresh rates failed due to: %v \n", err)
				}
			}
		}
	}()
}

const (
	appID = "872c2cd1f349476c8a93a24ea89f527c"
	url   = "https://openexchangerates.org/api/latest.json?app_id=%s"
)

type Result struct {
	Base  string             `json:"base"`
	Rates map[string]float32 `json:"rates"`
}

func refreshRates() error {
	// get rates from remote
	req, err := http.NewRequest("GET", fmt.Sprintf(url, appID), nil)
	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("fail to get rates from third-party, status code is %v", res.StatusCode))
	}
	body, _ := io.ReadAll(res.Body)

	var result Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	if len(result.Rates) == 0 {
		return errors.New("get 0 rates from third-party")
	}

	// create or update to db rate table
	dbClient, _ := db.NewClient()
	rateDB := db.NewRateDB(dbClient)

	for fromCurrency, fromRate := range result.Rates {
		for toCurrency, toRate := range result.Rates {
			rate := toRate / fromRate
			err = rateDB.AddRate(context.Background(), fromCurrency, toCurrency, rate)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
