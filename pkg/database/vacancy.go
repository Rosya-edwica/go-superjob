package database

import (
	"fmt"
	"log"
	"strings"

	"superjob/pkg/logger"
	"superjob/pkg/models"
)

// TODO: Insert Ignore

func (d *DB) SaveVacancy(v models.Vacancy) {
	query := `INSERT INTO vacancy (platform, id, url, name, city_id, position_id, prof_areas, specs, experience, salary_from, salary_to, key_skills, vacancy_date) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	db := d.GetDB()
	tx, _ := db.Begin()
	_, err := db.Exec(query, "superjob", v.Id, v.Url, v.Title, v.CityId, v.ProfessionId, v.ProfAreas, v.Specializations, v.Experience, v.SalaryFrom, v.SalaryTo, v.Skills, v.DateUpdate)
	checkErr(err)
	tx.Commit()
}

func (d *DB) SaveManyVacancies(vacancies []models.Vacancy) {
	query, valArgs := createQueryForMultipleInsertVacancies(vacancies)
	db := d.GetDB()
	tx, _ := db.Begin()
	_, err := db.Exec(query, valArgs...)
	if err != nil {
		log.Fatalf("Insert error: %s", err.Error())
	}
	tx.Commit()
	logger.Log.Printf("Сохранили %d вакансий", len(vacancies))

}

func createQueryForMultipleInsertVacancies(vacancies []models.Vacancy) (query string, valArgs []interface{}) {
	valStrings := []string{}
	valInsertCount := 1
	for _, v := range vacancies {
		valStrings = append(valStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", valInsertCount, valInsertCount+1, valInsertCount+2, valInsertCount+3, valInsertCount+4, valInsertCount+5, valInsertCount+6, valInsertCount+7, valInsertCount+8, valInsertCount+9, valInsertCount+10, valInsertCount+11, valInsertCount+12))
		valArgs = append(valArgs, "superjob", v.Id, v.Url, v.Title, v.CityId, v.ProfessionId, v.ProfAreas, v.Specializations, v.Experience, v.SalaryFrom, v.SalaryTo, v.Skills, v.DateUpdate)
		valInsertCount += 13
	}
	query = `INSERT INTO vacancy (platform, id, url, name, city_id, position_id, prof_areas, specs, experience, salary_from, salary_to, key_skills, vacancy_date) 
		VALUES` + strings.Join(valStrings, ",") + "ON CONFLICT DO NOTHING;"
	return
}
