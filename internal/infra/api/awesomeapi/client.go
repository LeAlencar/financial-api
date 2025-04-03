package awesomeapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lealencar/financial-api/internal/domain/models"
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

func (c *Client) GetExchangeRate(from, to string) (*models.Currency, error) {
	url := fmt.Sprintf("%s/%s-%s", c.BaseURL, from, to)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code inválido: %d", resp.StatusCode)
	}

	var result map[string]models.Currency
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	for _, currency := range result {
		return &currency, nil
	}

	return nil, fmt.Errorf("nenhum dado encontrado para o par %s-%s", from, to)
}
