package main

import (
	"fmt"
	"superjob/pkg/api"
	"superjob/pkg/database"
	"superjob/pkg/logger"
	"superjob/pkg/models"
	"sync"
)

var cities []models.City
var vacancies []models.Vacancy

func main() {
	db := database.InitDatabase()
	defer db.Close()

	cities = db.GetCities()
	positions := db.GetPositions()
	for _, position := range positions {
		findPositionVacancies(position)
		if len(vacancies) != 0 {
			db.SaveManyVacancies(vacancies)
			vacancies = []models.Vacancy{}
		}
		db.UpdatePositionStatusToParsed(position.Id)
	}

}

func findPositionVacancies(position models.Position) {
	names := position.OtherNames
	names = append(names, position.Name)
	for _, positionName := range names {
		if len(positionName) < 2 {
			continue
		}
		logger.Log.Printf("Ищем профессию: %s ", positionName)
		scrapePositionInRussia(positionName, position.Id)
	}
	return
}

func groupCities(count int) (groups [][]models.City) {
	for i := 0; i < len(cities); i += count {
		group := cities[i:]
		if len(group) >= 100 {
			group = group[:count]
		}
		groups = append(groups, group)
	}
	return
}

func scrapePositionInRussia(name string, id int) {
	superjob := api.API{PositionId: id, CityId: 0, CityEdwicaId: 0, PositionName: name}
	russiaUrl := superjob.CreateQuery()
	vacanciesCount := superjob.TotalVacanciesCount(russiaUrl)
	fmt.Println(name, vacanciesCount)
	if vacanciesCount > 500 {
		logger.Log.Printf("Тотечно ищем вакансии по всем городам. Количество: %d", vacanciesCount)
		scrapePosition(name, id)
	} else {
		if vacanciesCount != 0 {
			logger.Log.Printf("Ищем вакансии '%s' по всей России. Количество вакансий: %d", name, vacanciesCount)
			vacancies = append(vacancies, superjob.CollectAllVacancies(russiaUrl)...)
		} else {
			logger.Log.Printf("Ни найдено ни одной вакансии для: %s", name)
		}
	}
}

func scrapePosition(name string, id int) {
	cityGroups := groupCities(100)
	for _, group := range cityGroups {
		var wg sync.WaitGroup
		wg.Add(len(group))
		for _, city := range group {
			go scrapePositionInCity(name, id, city, &wg)
		}
		wg.Wait()
	}
	return
}

func scrapePositionInCity(name string, id int, city models.City, wg *sync.WaitGroup) {
	superjob := api.API{PositionId: id, CityId: city.SUPERJOB_ID, CityEdwicaId: city.EDWICA_ID, PositionName: name}
	baseUrl := superjob.CreateQuery()
	thisCityVacancies := superjob.CollectAllVacancies(baseUrl)
	logger.Log.Printf("В городе %s нашлось %d вакансий по запросу:%s", city.Name, len(thisCityVacancies), name)
	vacancies = append(vacancies, thisCityVacancies...)
	wg.Done()
}
