package handlers

import (
	"database/sql"

	models "forum/internal/database/models"
	"forum/internal/utils"
)

func GetCategories(db *sql.DB) ([]models.Category, error) {
	var result []models.Category
	var err error
	var rows *sql.Rows

	rows, err = utils.QueryRows(db, `SELECT id, name, description FROM categories`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var row models.Category
		if err := rows.Scan(&row.Id, &row.Name, &row.Description); err == nil {
			result = append(result, row)
		} else {
			return nil, err
		}
	}

	return result, nil
}
