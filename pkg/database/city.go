package database

import (
	"superjob/pkg/models"
)

func (d *DB) GetCities() (cities []models.City) {
	query := `SELECT id_edwica, name, id_superjob FROM h_city WHERE id_superjob != 0 ORDER BY id_superjob ASC`
	db := d.GetDB()
	rows, err := db.Query(query)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var name string
		var edwica_id, superjob_id int

		err = rows.Scan(&edwica_id, &name, &superjob_id)
		checkErr(err)
		cities = append(cities, models.City{
			EDWICA_ID:   edwica_id,
			SUPERJOB_ID: superjob_id,
			Name:        name,
		})
	}
	return
}
