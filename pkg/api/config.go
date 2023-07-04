package api

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"superjob/pkg/logger"
	"superjob/pkg/telegram"

	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

type API struct {
	PositionName string
	PositionId   int
	CityId       int
	CityEdwicaId int
}

const (
	TOKEN         = "v3.r.137188118.28a168bdc625692a956f581779284ebaed8f3141.06a12dccbe382b9dafaa8ea8cfbf11b87d65171e"
	CLIENT_SECRET = "v3.r.137188118.b1fc20f24580069b6f50e875ec950fda43eef2b5.ba7e377b764e0f2a8b6daff6e8078ba0cea471e6"
	CLIENT_ID     = "2093"
	MAIN_API_URL  = "https://api.superjob.ru/2.0/vacancies"
)

var HEADERS = map[string]string{
	"User-Agent":    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
	"Authorization": "Bearer " + TOKEN,
	"X-Api-App-Id":  CLIENT_SECRET,
}

func (api *API) CreateQuery() (query string) {
	params := url.Values{
		"count":             {"100"},
		"keywords[0][srws]": {"1"},              // Ищем в названии вакансии
		"keywords[0][skwc]": {"particular"},     // Ищем точную фразу
		"keywords[0][keys]": {api.PositionName}, // Фраза
	}
	if api.CityId != 0 {
		params.Add("town", strconv.Itoa(api.CityId))
	}

	return "https://api.superjob.ru/2.0/vacancies?" + params.Encode()
}

func getJson(url string, useHeaders bool) (string, error) {
	client := &http.Client{
		Timeout: 120 * time.Second,
	}
	request, err := http.NewRequest("GET", url, nil)
	checkErr(err)
	if useHeaders {
		for key, value := range HEADERS {
			request.Header.Set(key, value)
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if response.StatusCode != 200 {
		errMsg := gjson.Get(string(data), "error.message").String()
		return "", errors.New(errMsg)
	}

	return string(data), nil
}

func checkErr(err error) {
	if err != nil {
		telegram.Mailing(err.Error())
		logger.Log.Fatalln(err.Error())
	}
}
