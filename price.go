package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// PriceData represents ETH price information
type PriceData struct {
	Timestamp   int64              `json:"timestamp"`
	Eth         float64            `json:"eth_usd"`
	Sources     map[string]float64 `json:"sources"`
	Average     float64            `json:"average_usd"`
	Change24h   float64            `json:"change_24h_percent"`
}

// fetchETHPrice fetches ETH/USD price from multiple sources
func fetchETHPrice() (*PriceData, error) {
	sources := make(map[string]float64)
	
	// Try CoinGecko
	if price, err := fetchCoinGeckoPrice(); err == nil {
		sources["coingecko"] = price
	}
	
	// Try Coinbase
	if price, err := fetchCoinbasePrice(); err == nil {
		sources["coinbase"] = price
	}
	
	// Try Kraken
	if price, err := fetchKrakenPrice(); err == nil {
		sources["kraken"] = price
	}
	
	if len(sources) == 0 {
		return nil, fmt.Errorf("failed to fetch price from all sources")
	}
	
	// Calculate average
	var sum float64
	for _, price := range sources {
		sum += price
	}
	average := sum / float64(len(sources))
	
	return &PriceData{
		Timestamp: time.Now().Unix(),
		Eth:       round(average, 2),
		Sources:   sources,
		Average:   round(average, 2),
		Change24h: 0, // Would need historical data
	}, nil
}

func fetchCoinGeckoPrice() (float64, error) {
	resp, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	
	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	
	if eth, ok := result["ethereum"]; ok {
		if price, ok := eth["usd"]; ok {
			return price, nil
		}
	}
	return 0, fmt.Errorf("price not found")
}

func fetchCoinbasePrice() (float64, error) {
	resp, err := http.Get("https://api.coinbase.com/v2/exchange-rates?currency=ETH")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	
	var result struct {
		Data struct {
			Rates map[string]string `json:"rates"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	
	if rate, ok := result.Data.Rates["USD"]; ok {
		return strconv.ParseFloat(rate, 64)
	}
	return 0, fmt.Errorf("rate not found")
}

func fetchKrakenPrice() (float64, error) {
	resp, err := http.Get("https://api.kraken.com/0/public/Ticker?pair=ETHUSD")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	
	var result struct {
		Result map[string]struct {
			C []string `json:"c"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	
	for _, v := range result.Result {
		if len(v.C) > 0 {
			return strconv.ParseFloat(v.C[0], 64)
		}
	}
	return 0, fmt.Errorf("ticker not found")
}
