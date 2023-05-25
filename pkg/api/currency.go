package api

import (
	"fmt"
	"log"
	"superjob/pkg/models"

	"github.com/tidwall/gjson"
)

var RelevantCurrencies = []models.Currency{}

func init() {
	RelevantCurrencies = GetCurrencies()
}

func GetCurrencies() (currencies []models.Currency) {
	json, err := getJson("https://api.hh.ru/dictionaries", false)
	fmt.Println("Обновляем валюту!")
	if err != nil {
		log.Printf("Не удалось обновить валюту. Текст сообщения: %s", err)
	}
	for _, item := range gjson.Get(json, "currency").Array() {
		var abbr string
		switch item.Get("abbr").String() {
		case "руб.":
			abbr = "rub"
		case "грн.":
			abbr = "uah"
		case "сум":
			abbr = "uzs"
		}
		currencies = append(currencies, models.Currency{
			Abbr: abbr,
			Name: item.Get("name").String(),
			Rate: item.Get("rate").Float(),
		})
	}
	return
}
