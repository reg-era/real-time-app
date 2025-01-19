package handlers

import (
	"database/sql"
	"net/http"

	"forum/internal/database"
	models "forum/internal/database/models"
	"forum/internal/utils"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, userId int) {
	if r.URL.Query().Has("category") {
		category := r.URL.Query().Get("category")
		result, err := utils.QueryRow(db, `SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?)`, category)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
			return
		}

		var exists bool
		if err := result.Scan(&exists); err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
			return
		}
		if !exists {
			utils.RespondWithJSON(w, http.StatusNotFound, utils.ErrorResponse{Error: "Status Not Found"})
			return
		}

		postIds, err := database.GetCategoryContentIds(db, category)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, utils.ErrorResponse{Error: "Internal Server Error"})
			return
		}

		data := struct {
			PostIds []int `json:"post_ids"`
		}{PostIds: postIds}

		utils.RespondWithJSON(w, http.StatusOK, data)
	}
}

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
