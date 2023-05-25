package database

import (
	"superjob/pkg/models"
)

func (d *DB) GetCities() (cities []models.City) {
	query := `SELECT id_hh, id_edwica, id_rabota_ru, name, id_superjob FROM city WHERE id_superjob != 0  and id_hh != 0 ORDER BY id_hh ASC`
	db := d.GetDB()
	rows, err := db.Query(query)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var name, rabota_id string
		var hh_id, edwica_id, superjob_id int

		err = rows.Scan(&hh_id, &edwica_id, &rabota_id, &name, &superjob_id)
		checkErr(err)
		cities = append(cities, models.City{
			HH_ID:       hh_id,
			EDWICA_ID:   edwica_id,
			SUPERJOB_ID: superjob_id,
			RABOTA_ID:   rabota_id,
			Name:        name,
		})
	}
	return
}
