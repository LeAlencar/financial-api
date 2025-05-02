package main

import (
	"testing"
	"time"

	"github.com/leandroalencar/banco-dados/shared/models"
)

func TestConvertAPIToModel(t *testing.T) {
	api := models.QuotationAPI{
		Code:      "USD",
		Codein:    "BRL",
		Bid:       "5.1234",
		Ask:       "5.2345",
		Timestamp: "1555360543",
	}
	expected := models.Quotation{
		CurrencyPair:  "USD/BRL",
		BuyPrice:      5.1234,
		SellPrice:     5.2345,
		Timestamp:     time.Unix(1555360543, 0),
		LastUpdatedBy: "awesomeapi",
	}
	result, err := ConvertAPIToModel(api)
	if err != nil {
		t.Fatalf("Erro ao converter: %v", err)
	}
	if result.CurrencyPair != expected.CurrencyPair {
		t.Errorf("CurrencyPair = %v, esperado %v", result.CurrencyPair, expected.CurrencyPair)
	}
	if result.BuyPrice != expected.BuyPrice {
		t.Errorf("BuyPrice = %v, esperado %v", result.BuyPrice, expected.BuyPrice)
	}
	if result.SellPrice != expected.SellPrice {
		t.Errorf("SellPrice = %v, esperado %v", result.SellPrice, expected.SellPrice)
	}
	if !result.Timestamp.Equal(expected.Timestamp) {
		t.Errorf("Timestamp = %v, esperado %v", result.Timestamp, expected.Timestamp)
	}
	if result.LastUpdatedBy != expected.LastUpdatedBy {
		t.Errorf("LastUpdatedBy = %v, esperado %v", result.LastUpdatedBy, expected.LastUpdatedBy)
	}
}
