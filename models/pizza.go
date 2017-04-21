package models

import (
	"database/sql"
	"log"

	"fmt"

	"github.com/lempiy/pizza-app-pq/types"
	"github.com/lempiy/pizza-app-pq/utils/utils"
)

// GetAllPizzas retrieves all the pizzas
func GetAllPizzas() ([]types.Pizza, error) {
	log.Println("All getting pizzas")
	pizzas := make([]types.Pizza, 0)
	var pizza types.Pizza
	var rows *sql.Rows

	querySQL := `SELECT p.id, p.name, p.size, p.description, p.img_url, p.accepted, p.created_date,
                  p.updated_date, p.deleted, p.price, c.id, c.name, c.description, u.id, u.username, u.created_date
                      FROM pizza p
                          LEFT JOIN pizza_category c ON c.id = p.category_id
                          LEFT JOIN person u ON u.id = p.user_id
                      WHERE p.deleted = 0 AND p.accepted = 0
                      ORDER BY p.created_date ASC`

	rows = database.query(querySQL)

	if err != nil {
		log.Println("Failed to get all pizzas from db")
	}

	defer rows.Close()
	for rows.Next() {
		pizza = types.Pizza{}
		var deleted int
		var accepted int
		err = rows.Scan(
			&pizza.ID, &pizza.Name, &pizza.Size, &pizza.Description,
			&pizza.ImgURL, &accepted, &pizza.Created, &pizza.Updated, &deleted, &pizza.Price,
			&pizza.Category.ID, &pizza.Category.Name, &pizza.Category.Description,
			&pizza.Author.ID, &pizza.Author.Name, &pizza.Author.Created)

		pizza.Deleted = utils.Itob(deleted)
		pizza.Accepted = utils.Itob(accepted)

		ingridients, err := GetIngridientsOfPizza(pizza.ID)

		if err != nil {
			log.Printf("Failed to get ingredients for pizza with ID %d", pizza.ID)
		}

		pizza.Ingredients = ingridients

		pizzas = append(pizzas, pizza)
	}
	return pizzas, nil
}

// GetPizzaByID retrieves some pizza by ID, returns pizza and is_my_pizza bool
func GetPizzaByID(id int) (types.Pizza, error) {
	log.Println("All getting pizzas")
	var pizza types.Pizza
	var rows *sql.Rows

	querySQL := `SELECT p.id, p.name, p.size, p.description, p.img_url, p.accepted, p.created_date,
                  p.updated_date, p.deleted, p.price, c.id, c.name, c.description, u.id, u.username, u.created_date
                      FROM pizza p
                          LEFT JOIN pizza_category c ON c.id = p.category_id
                          LEFT JOIN person u ON u.id = p.user_id
                      WHERE p.id = $1 AND p.deleted = 0 AND p.accepted = 0;`

	rows = database.query(querySQL, id)

	if err != nil {
		log.Println("Failed to get pizza by id")
	}

	defer rows.Close()
	for rows.Next() {
		var deleted int
		var accepted int
		err = rows.Scan(
			&pizza.ID, &pizza.Name, &pizza.Size, &pizza.Description,
			&pizza.ImgURL, &accepted, &pizza.Created, &pizza.Updated, &deleted, &pizza.Price,
			&pizza.Category.ID, &pizza.Category.Name, &pizza.Category.Description,
			&pizza.Author.ID, &pizza.Author.Name, &pizza.Author.Created)

		pizza.Deleted = utils.Itob(deleted)
		pizza.Accepted = utils.Itob(accepted)

		ingridients, err := GetIngridientsOfPizza(pizza.ID)

		if err != nil {
			log.Printf("Failed to get ingredients for pizza with ID %d", pizza.ID)
		}

		pizza.Ingredients = ingridients
	}
	return pizza, nil
}

//PostPizza creates pizza order
func PostPizza(pizza types.PizzaPost, price float64) error {
	sqlQuery := `INSERT INTO pizza(name,user_id,category_id,size,description,img_url,price,created_date,updated_date)
                VALUES($1, $2, $3, $4, $5, $6, $7,now(), now())`
	lastID, err := insertWithReturningID(sqlQuery, pizza.Name, pizza.AuthorID, pizza.CategoryID,
		pizza.Size, pizza.Description, pizza.ImgURL, price)
	err = PostUsedIngredients(lastID, pizza.IngredientsIDs)
	return err
}

//PostPizzaParallel creates pizza order async
func PostPizzaParallel(pizza types.PizzaPost, price float64, done chan bool) {
	err := PostPizza(pizza, price)
	if err != nil {
		done <- false
		log.Println(err)
	} else {
		done <- true
	}
}

//PostUsedIngredients - places all ingr used in particular pizza
func PostUsedIngredients(pizzaID int, ingrIDs []int) error {
	sqlQuery := `INSERT INTO used_ingredient(ingredient_id, pizza_id) VALUES`
	var err error
	for i, id := range ingrIDs {
		var queryValue string
		if i == 0 {
			queryValue = fmt.Sprintf("(%d, %d)", id, pizzaID)
		} else {
			queryValue = fmt.Sprintf(", (%d, %d)", id, pizzaID)
		}
		sqlQuery += queryValue
		if err != nil {
			return err
		}
	}
	sqlQuery += `;`
	log.Println(sqlQuery)
	err = singleQuery(sqlQuery)
	return err
}

//CalculatePizzaPrice - calculates total price of some pizza
func CalculatePizzaPrice(ingredients []int, size int) (float64, error) {
	ingreds, err := GetIngredientsByIDs(ingredients)
	var totalPrice float64
	totalPrice = 10.00
	if err != nil {
		log.Println("Error upon price calculation.")
		return totalPrice, err
	}
	for _, ingr := range ingreds {
		totalPrice += float64(ingr.Price)
	}

	totalPrice = float64(size) * totalPrice / float64(60)

	return totalPrice, err
}

//DeletePizza sets deleted flag to pizza
func DeletePizza(id int, userID int) (int, error) {
	sqlQuery := `UPDATE pizza
                  SET deleted = 1
                WHERE id=$1 AND user_id=$2;`
	deletedCount, err := singleQueryWithAffected(sqlQuery, id, userID)
	if err != nil {
		panic(err)
	}
	return deletedCount, err
}

//RemoveUnusedIngredients deletes unused ingredients from used_ingredient table
func RemoveUnusedIngredients(pizzaID int) (int, error) {
	sqlQuery := `DELETE FROM used_ingredient
                WHERE pizza_id=$1;`
	deletedCount, err := singleQueryWithAffected(sqlQuery, pizzaID)
	if err != nil {
		panic(err)
	}
	return deletedCount, err
}

//UpdatePizza applies change to pizza by its ID
func UpdatePizza(pizza types.PizzaPost, price float64, id int, userID int) error {
	sqlQuery := `UPDATE pizza
                  SET name = $1,user_id = $2,category_id = $3, size = $4,
                      description = $5, img_url = $6, price = $7,
                      updated_date = now()
                WHERE id=$8 AND user_id=$9;`
	log.Printf("pizza %v | price %v | id %v | userID %v\n", pizza, price, id, userID)
	err := singleQuery(sqlQuery, pizza.Name, pizza.AuthorID, pizza.CategoryID,
		pizza.Size, pizza.Description, pizza.ImgURL, price, id, userID)

	err = PostUsedIngredients(id, pizza.IngredientsIDs)
	return err
}

//UpdatePizzaParallel upadtes pizza order async
func UpdatePizzaParallel(pizza types.PizzaPost, price float64,
	id int, userID int, done chan bool) {
	err := UpdatePizza(pizza, price, id, userID)
	if err != nil {
		done <- false
		log.Println(err)
	} else {
		done <- true
	}
}

//AcceptPizza sets deleted flag to pizza
func AcceptPizza(id int, userID int) (int, error) {
	sqlQuery := `UPDATE pizza
                  SET accepted = 1
                WHERE id=$1 AND user_id=$2;`
	acceptedCount, err := singleQueryWithAffected(sqlQuery, id, userID)
	if err != nil {
		panic(err)
	}
	return acceptedCount, err
}
