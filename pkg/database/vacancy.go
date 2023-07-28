package database

import (
	"log"

	"superjob/pkg/logger"
	"superjob/pkg/models"
)


func (d *DB) SaveManyVacancies(vacancies []models.Vacancy) {
	query, valArgs := createQueryForMultipleInsertVacanciesMYSQL(vacancies)
	db := d.GetDB()
	tx, _ := db.Begin()
	_, err := db.Exec(query, valArgs...)
	if err != nil {
		log.Fatalf("Insert error: %s", err.Error())
	}
	tx.Commit()
	logger.Log.Printf("Сохранили %d вакансий", len(vacancies))

}


func createQueryForMultipleInsertVacanciesMYSQL(vacancies []models.Vacancy) (query string, valArgs []interface{}) {
	query = "INSERT IGNORE INTO h_vacancy (id, name, url, city_id, position_id, prof_areas, specs, experience, salary_from, salary_to, key_skills, vacancy_date, platform) VALUES "
	for _, v := range vacancies {
		query += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),"
		valArgs = append(valArgs,  v.Id, v.Title, v.Url, v.CityId, v.ProfessionId, v.ProfAreas, v.Specializations, v.Experience, v.SalaryFrom, v.SalaryTo, v.Skills, v.DateUpdate, "superjob")
	}
	query = query[0:len(query)-1]
	return
}