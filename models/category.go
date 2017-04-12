package models

import (
	"database/sql"
	"log"

	"github.com/lempiy/pizza-app-pq/types"
	"github.com/lempiy/pizza-app-pq/utils/utils"
)

// GetAllCategories retrieves all the categories
func GetAllCategories() ([]types.Category, error) {
	log.Println("All getting categories")
	categories := make([]types.Category, 0)
	var categoryDB types.CategoryDB
	var rows *sql.Rows

	querySQL := `SELECT * FROM pizza_category;`

	rows = database.query(querySQL)

	if err != nil {
		log.Println("Failed to get all categories from db")
	}

	defer rows.Close()
	for rows.Next() {
		categoryDB = types.CategoryDB{}

		err = rows.Scan(
			&categoryDB.ID, &categoryDB.Name, &categoryDB.Description, &categoryDB.IsDefault)

		category := types.Category{
			ID:          categoryDB.ID,
			Name:        categoryDB.Name,
			Description: categoryDB.Description,
			IsDefault:   utils.Itob(categoryDB.IsDefault)}

		categories = append(categories, category)
	}
	return categories, nil
}
