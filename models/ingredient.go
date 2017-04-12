package models

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/lempiy/pizza-app-pq/types"
)

// GetIngridientsOfPizza gets all ingridients on some pizza
func GetIngridientsOfPizza(pizzaID int) ([]types.IngredientSimple, error) {
	log.Printf("Get all ingridients for pizza %d\n", pizzaID)
	var ingredients []types.IngredientSimple
	var ingredient types.IngredientSimple
	var rows *sql.Rows

	querySQL := `SELECT i.id, i.name, i.image_url, i.price
                FROM used_ingredient n
                  LEFT JOIN ingredient i ON i.id=n.ingredient_id
                WHERE n.pizza_id=$1 ;`

	rows = database.query(querySQL, pizzaID)

	defer rows.Close()
	for rows.Next() {
		ingredient = types.IngredientSimple{}

		err = rows.Scan(&ingredient.ID, &ingredient.Name, &ingredient.ImgURL, &ingredient.Price)
		if err != nil {
			log.Println(err)
		}

		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}

//GetAllIngredients used to get all ingredients from DB
func GetAllIngredients() ([]types.Ingredient, error) {
	log.Println("All getting categories")
	ingredients := make([]types.Ingredient, 0)
	var ingredient types.Ingredient
	var rows *sql.Rows

	querySQL := `SELECT i.id, i.name, i.description, i.image_url, i.price, i.created_date, u.id, u.username, u.created_date
                FROM ingredient i
                  JOIN person u ON u.id=i.user_id;`

	rows = database.query(querySQL)

	if err != nil {
		log.Println("Failed to get all categories from db")
	}

	defer rows.Close()
	for rows.Next() {
		ingredient = types.Ingredient{}

		err = rows.Scan(
			&ingredient.ID, &ingredient.Name, &ingredient.Description,
			&ingredient.ImgURL, &ingredient.Price, &ingredient.Created,
			&ingredient.Author.ID, &ingredient.Author.Name, &ingredient.Author.Created)

		ingredients = append(ingredients, ingredient)
	}
	return ingredients, nil
}

//GetIngredientsByIDs - gets a bunch of ingredients by ids
func GetIngredientsByIDs(ingrdIDs []int) ([]types.Ingredient, error) {
	log.Println("All getting categories")
	var ingredients []types.Ingredient
	var ingredient types.Ingredient
	var rows *sql.Rows

	querySQL := `SELECT i.id, i.name, i.price
                FROM ingredient i
                WHERE id in `

	queryIDs := strings.Trim(strings.Replace(fmt.Sprint(ingrdIDs), " ", ",", -1), "[]")
	querySQL += "(" + queryIDs + ");"

	rows = database.query(querySQL)

	if err != nil {
		log.Println("Failed to get all categories from db")
	}

	defer rows.Close()
	for rows.Next() {
		ingredient = types.Ingredient{}

		err = rows.Scan(
			&ingredient.ID, &ingredient.Name, &ingredient.Price)

		ingredients = append(ingredients, ingredient)
	}
	return ingredients, nil
}
