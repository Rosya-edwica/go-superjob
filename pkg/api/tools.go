package api

import (
	"time"

	"github.com/tidwall/gjson"
)

type Salary struct {
	From     float64
	To       float64
	Currency string
}

func getProfAreas(vacancyJson gjson.Result) (profAreas []string) {

	for _, item := range vacancyJson.Get("catalogues").Array() {
		profAreas = append(profAreas, item.Get("title").String())
	}
	return
}
func getSpecs(vacancyJson gjson.Result) (specs []string) {
	for _, profArea := range vacancyJson.Get("catalogues").Array() {
		for _, item := range profArea.Get("positions").Array() {
			specs = append(specs, item.Get("title").String())
		}
	}
	return
}

func convertDateUpdate(timestamp int64) (date string) {
	dateUpdate := time.Unix(timestamp, 0)
	date = dateUpdate.Format(time.RFC3339Nano)
	return
}

func convertExperienceId(id int) (experience string) {
	switch id {
	case 2:
		experience = "От 1 года до 3 лет"
	case 3:
		experience = "От 3 до 6 лет"
	case 4:
		experience = "От 6 лет"
	default:
		experience = "Нет опыта"
	}
	return
}

func getSalary(jsonVacancy gjson.Result) (salary Salary) {
	salary.From = jsonVacancy.Get("payment_from").Float()
	salary.To = jsonVacancy.Get("payment_to").Float()
	salary.Currency = jsonVacancy.Get("currency").String()

	switch salary.Currency {
	case "rub":
		return salary
	case "":
		return Salary{}
	default:
		return convertSalaryToRUR(salary)
	}

}

func convertSalaryToRUR(salary Salary) (convertedSalary Salary) {
	for _, cur := range RelevantCurrencies {
		if cur.Abbr == salary.Currency {
			salary.To = salary.To / cur.Rate
			salary.From = salary.From / cur.Rate
			salary.Currency = "rub"
			return salary
		}
	}

	return
}
