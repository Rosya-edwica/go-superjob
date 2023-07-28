package database

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"superjob/pkg/models"
)

func (d *DB) GetProfessionsFromFile(fromFile bool) (positions []models.Position) {
	var query string
	if fromFile {
		areas := strings.ToLower(arrayToPostgresList(readProfAreasFromFile()))
		query = `SELECT position.id, position.name, position.other_names
		FROM position
		LEFT JOIN position_to_prof_area ON position_to_prof_area.position_id=position.id
		LEFT JOIN prof_area_to_specialty ON prof_area_to_specialty.id=position_to_prof_area.area_id
		LEFT JOIN professional_area ON professional_area.id=prof_area_to_specialty.prof_area_id
		WHERE LOWER(professional_area.name) IN ` + areas
		fmt.Println("Парсим профессии этих профобластей: ", areas)
	} else {
		query = "SELECT id, name, other_names FROM positon"
		fmt.Println("Парсим абсолютно все профессии")
	}

	db := d.GetDB()
	rows, err := db.Query(query)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var (
			name string
			other sql.NullString
			id int
		)
		err = rows.Scan(&id, &name, &other)
		checkErr(err)

		prof := models.Position{
			Id:         id,
			Name:       name,
			OtherNames: strings.Split(other.String, "|"),
		}
		positions = append(positions, prof)

	}
	return
}


func readProfAreasFromFile() (areas []string) {
	filepath := "prof_areas.txt"
	file, err := os.Open(filepath)
	if err != nil {
		panic("Создайте файл prof_areas.txt для парсинга по профобластям")
	}
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		areas = append(areas, fileScanner.Text())
	}
	file.Close()
	return
}

func arrayToPostgresList(items []string) (result string) {
	var updatedList []string
	for _, i := range items {
		updatedList = append(updatedList, fmt.Sprintf("'%s'", i))
	}
	result = "(" + strings.Join(updatedList, ",") + ")"
	return
}
