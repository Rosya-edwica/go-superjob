package database

import (
	"fmt"
	"strings"

	"superjob/pkg/logger"
	"superjob/pkg/models"
)

func (d *DB) GetPositions() (positions []models.Position) {
	query := `SELECT id, name, other_names FROM position WHERE parsed = false AND id != 0`
	db := d.GetDB()
	rows, err := db.Query(query)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var name, other_names string
		var id int

		err = rows.Scan(&id, &name, &other_names)
		checkErr(err)
		positions = append(positions, models.Position{
			Id:         id,
			Name:       name,
			OtherNames: strings.Split(other_names, "|"),
		})
	}
	return
}

func (d *DB) UpdatePositionStatusToParsed(positionId int) {
	query := fmt.Sprintf("UPDATE position SET parsed=true WHERE id = %d", positionId)
	db := d.GetDB()
	tx, _ := db.Begin()
	_, err := db.Exec(query)
	checkErr(err)
	tx.Commit()
	logger.Log.Printf("Спарсили профессию: %d", positionId)
}
