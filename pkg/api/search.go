package api

import (
	"fmt"
	"strings"
	"superjob/pkg/database"
	"superjob/pkg/logger"
	"superjob/pkg/models"
	"time"

	"github.com/tidwall/gjson"
)

var AllCities []models.City

func init() {
	db := database.InitDatabase()
	AllCities = db.GetCities()
}

func (api *API) TotalVacanciesCount(url string) (count int64) {
	json, err := getJson(url, true)
	checkErr(err)
	count = gjson.Get(json, "total").Int()
	return
}

func (api *API) CollectAllVacancies(url string) (vacancies []models.Vacancy) {
	var pageNum int
	for {
		pageVacancies := api.GetVacanciesFromPage(fmt.Sprintf("%s&page=%d", url, pageNum))
		if len(pageVacancies) == 0 {
			break
		}
		vacancies = append(vacancies, pageVacancies...)
		pageNum++
	}
	return
}

func (api *API) GetVacanciesFromPage(url string) (vacancies []models.Vacancy) {
	logger.Log.Println(url)
	json, err := getJson(url, true)
	if err != nil {
		if err.Error() == "Limit is over" {
			time.Sleep(time.Second * 3)
			json, err = getJson(url, true)
		} else {
			logger.Log.Printf("ОШИБКА: %s", err.Error())
			return
		}
	}

	for _, item := range gjson.Get(json, "objects").Array() {
		var vacancy models.Vacancy
		salary := getSalary(item)
		
		if api.CityEdwicaId != 0 {
			vacancy.CityId = api.CityEdwicaId
		} else {
			vacancy.CityId = setEdwicaCityId(item)
		}
		vacancy.Id = int(item.Get("id").Int())
		vacancy.ProfessionId = api.PositionId
		vacancy.Title = item.Get("profession").String()
		vacancy.SalaryFrom = salary.From
		vacancy.SalaryTo = salary.To
		vacancy.Url = item.Get("link").String()
		vacancy.ProfAreas = strings.Join(getProfAreas(item), "|")
		vacancy.Specializations = strings.Join(getSpecs(item), "|")
		vacancy.Experience = convertExperienceId(int(item.Get("experience.id").Int()))
		vacancy.DateUpdate = convertDateUpdate(item.Get("date_published").Int())

		vacancies = append(vacancies, vacancy)
	}
	return
}


func setEdwicaCityId(json gjson.Result) (id int) {
	superjobId := int(json.Get("town.id").Int())
	for _, city := range AllCities {
		if city.SUPERJOB_ID == superjobId {
			return city.EDWICA_ID
		}
	}
	fmt.Println("Не удалось подобрать город")
	return
}