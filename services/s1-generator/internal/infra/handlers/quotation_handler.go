package handlers

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leandroalencar/banco-dados/shared/models"
	"github.com/leandroalencar/banco-dados/shared/utils"
)

type QuotationHandler struct {
	rabbitmq *utils.RabbitMQ
}

func NewQuotationHandler(rabbitmq *utils.RabbitMQ) *QuotationHandler {
	return &QuotationHandler{
		rabbitmq: rabbitmq,
	}
}

// GenerateQuotations fetches current USD-BRL quotation and generates multiple variations
func (h *QuotationHandler) GenerateQuotations(c *gin.Context) {
	// Get the number of quotations to generate (default to 100)
	count := 100
	if countStr := c.Query("count"); countStr != "" {
		if n, err := strconv.Atoi(countStr); err == nil && n > 0 {
			count = n
		}
	}

	// Fetch base quotation from API
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quotations"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read API response"})
		return
	}

	// Parse API response
	var apiQuotations map[string]models.QuotationAPI
	if err := json.Unmarshal(body, &apiQuotations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API response"})
		return
	}

	// Get the base USD-BRL quotation
	baseQuotation, ok := apiQuotations["USDBRL"]
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "USD-BRL quotation not found in response"})
		return
	}

	// Convert base quotation
	baseQuote, err := ConvertAPIToModel(baseQuotation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert base quotation"})
		return
	}

	var publishedQuotations []models.Quotation
	baseTime := time.Now().Add(-time.Duration(count) * time.Minute) // Start from count minutes ago

	// Generate and publish multiple quotations
	for i := 0; i < count; i++ {
		// Create a variation of the base quotation
		quotation := generateQuotationVariation(baseQuote, baseTime.Add(time.Duration(i)*time.Minute))

		// Publish to RabbitMQ
		err = h.rabbitmq.PublishMessage("quotations", quotation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish quotation to queue"})
			return
		}

		publishedQuotations = append(publishedQuotations, quotation)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quotations generated and published successfully",
		"count":   len(publishedQuotations),
		"data":    publishedQuotations,
	})
}

// generateQuotationVariation creates a variation of the base quotation
func generateQuotationVariation(base models.Quotation, timestamp time.Time) models.Quotation {
	// Generate random variations within 0.5% of the base price
	variation := 0.005 // 0.5%

	// Random multiplier between -1 and 1
	buyMultiplier := (rand.Float64()*2 - 1)
	sellMultiplier := (rand.Float64()*2 - 1)

	// Calculate new prices with variations
	buyVariation := base.BuyPrice * variation * buyMultiplier
	sellVariation := base.SellPrice * variation * sellMultiplier

	return models.Quotation{
		CurrencyPair:  base.CurrencyPair,
		BuyPrice:      base.BuyPrice + buyVariation,
		SellPrice:     base.SellPrice + sellVariation,
		Timestamp:     timestamp,
		LastUpdatedBy: "awesomeapi-generator",
	}
}

// ConvertAPIToModel converts API response to Quotation model
func ConvertAPIToModel(api models.QuotationAPI) (models.Quotation, error) {
	buyPrice, err := strconv.ParseFloat(api.Bid, 64)
	if err != nil {
		return models.Quotation{}, err
	}

	sellPrice, err := strconv.ParseFloat(api.Ask, 64)
	if err != nil {
		return models.Quotation{}, err
	}

	tsInt, err := strconv.ParseInt(api.Timestamp, 10, 64)
	if err != nil {
		return models.Quotation{}, err
	}
	timestamp := time.Unix(tsInt, 0)

	return models.Quotation{
		CurrencyPair:  api.Code + "/" + api.Codein,
		BuyPrice:      buyPrice,
		SellPrice:     sellPrice,
		Timestamp:     timestamp,
		LastUpdatedBy: "awesomeapi",
	}, nil
}
