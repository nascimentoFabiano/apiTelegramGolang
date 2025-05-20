package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	telebot "gopkg.in/telebot.v4"
)

type CurrencyResponseDolar struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func getMoedaUSD() (string, error) {

	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	// Fazendo a requisição GET
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	// Lendo o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %w", err)
	}

	// Parse do JSON
	var data CurrencyResponseDolar
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("erro ao parsear JSON: %w", err)
	}

	return data.USDBRL.High, nil

}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar variáveis de ambiente: ", err)
		panic(err)
	}

	secretKeyFatherApi := os.Getenv("SECRET_KEY_API_FATHER_TELEGRAM")

	pref := telebot.Settings{
		Token:  secretKeyFatherApi,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	markup := &telebot.ReplyMarkup{}
	btnDolar := markup.Data("💵 Dólar", "dolar", "cotacao_dolar")
	btnEuro := markup.Data("💶 Euro", "euro", "cotacao_euro")

	markup.Inline(
		markup.Row(btnDolar, btnEuro),
	)

	b.Handle("/start", func(c telebot.Context) error {
		return c.Send("Oooooo meu amiguinho, qual caô de hoje irmao?")
	})

	b.Handle("/cotacoes", func(c telebot.Context) error {
		return c.Send("qual moeda quer consultar?", markup)
	})

	b.Handle(&btnDolar, func(c telebot.Context) error {
		high, err := getMoedaUSD()
		if err != nil {
			return c.Send("Erro ao obter cotação do DOLAR: ", err)
		}

		return c.Send(fmt.Sprintf("Cotação do Dólar: R$ %s", high))
	})

	// b.Handle(&btnEuro, func(c telebot.Context) error {
	// 	// cotacoes, err := buscarCotacaoAwesomeAPI()
	// 	// if err != nil {
	// 	// 	// Se der erro, exibe a mensagem de erro
	// 	// 	fmt.Println("Erro ao obter cotação do Euro:", err)
	// 	// 	return c.Send("Erro ao obter cotação do Euro.")
	// 	// }
	// 	// // Exibe a cotação do Euro (Bid)
	// 	// return c.Send(fmt.Sprintf("Cotação do Euro: R$ %s", cotacoes.EuroBRL.Bid))
	// })

	log.Print("Aplicação inciando!")
	b.Start()

}
