package awesomeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/leandroalencar/banco-dados/shared/models"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL:    "https://economia.awesomeapi.com.br/json/last",
		HTTPClient: &http.Client{},
	}
}

func (c *Client) GetExchangeRate(from, to string) (*models.Quotation, error) {
	url := fmt.Sprintf("%s/%s-%s", c.BaseURL, from, to)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error in request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	var result map[string]struct {
		Bid       string `json:"bid"`
		Ask       string `json:"ask"`
		Timestamp string `json:"timestamp"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	for _, rate := range result {
		bid, _ := strconv.ParseFloat(rate.Bid, 64)
		ask, _ := strconv.ParseFloat(rate.Ask, 64)
		timestamp, _ := strconv.ParseInt(rate.Timestamp, 10, 64)

		return &models.Quotation{
			CurrencyPair:  fmt.Sprintf("%s/%s", from, to),
			BuyPrice:      bid,
			SellPrice:     ask,
			Timestamp:     time.Unix(timestamp, 0),
			LastUpdatedBy: "awesomeapi",
		}, nil
	}

	return nil, fmt.Errorf("no data found for pair %s-%s", from, to)
}
