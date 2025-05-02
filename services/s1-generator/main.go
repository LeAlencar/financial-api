package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/leandroalencar/banco-dados/shared/models"
	"github.com/leandroalencar/banco-dados/shared/utils"
)

// QuotationAPI representa o formato da resposta da API
type QuotationAPI struct {
	Code      string `json:"code"`
	Codein    string `json:"codein"`
	Bid       string `json:"bid"`
	Ask       string `json:"ask"`
	Timestamp string `json:"timestamp"`
}

// ConvertAPIToModel converte a resposta da API para o model Quotation
func ConvertAPIToModel(api QuotationAPI) (models.Quotation, error) {
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

func main() {
	// Conectar ao RabbitMQ
	conn, ch, err := utils.ConnectRabbitMQ()
	if err != nil {
		log.Fatalf("Erro ao conectar no RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	queueName := "cotacoes"

	for {
		// Buscar cotações da API
		url := "https://economia.awesomeapi.com.br/json/daily/USD-BRL/4"
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Erro ao buscar API: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("Erro ao ler resposta da API: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		// Deserializar para slice de QuotationAPI
		var apiCotacoes []QuotationAPI
		if err := json.Unmarshal(body, &apiCotacoes); err != nil {
			log.Printf("Erro ao deserializar JSON: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		// Converter e publicar cada cotação
		for _, apiCotacao := range apiCotacoes {
			cotacao, err := ConvertAPIToModel(apiCotacao)
			if err != nil {
				log.Printf("Erro ao converter cotação: %v", err)
				continue
			}
			msg, _ := json.Marshal(cotacao)
			err = utils.PublishMessage(ch, queueName, msg)
			if err != nil {
				log.Printf("Erro ao publicar no RabbitMQ: %v", err)
			} else {
				log.Printf("Cotação publicada: %s - Compra: %.2f, Venda: %.2f", cotacao.CurrencyPair, cotacao.BuyPrice, cotacao.SellPrice)
			}
		}

		// Espera antes de buscar novamente
		time.Sleep(30 * time.Second)
	}
}
