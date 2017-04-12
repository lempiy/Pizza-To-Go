package types

import (
	"github.com/dgrijalva/jwt-go"
)

// UserClaims for jwt token
type UserClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// User type used for pizza authors
type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Created string `json:"created"`
}

// Pizzas collection type
type Pizzas []Pizza

// Pizza stores pizza values
type Pizza struct {
	ID          int                `json:"id"`
	Name        string             `json:"name"`
	Author      User               `json:"author"`
	Category    Category           `json:"category"`
	Size        int                `json:"size"`
	Description string             `json:"description"`
	Created     string             `json:"created_date"`
	Updated     string             `json:"updated_date"`
	ImgURL      string             `json:"img_url"`
	Ingredients []IngredientSimple `json:"ingredients"`
	Deleted     bool               `json:"deleted"`
	Accepted    bool               `json:"accepted, omitempty"`
	Price       float64            `json:"price, omitempty"`
}

// PizzaPost - type of user pizza data
type PizzaPost struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	AuthorID       int     `json:"author_id"`
	CategoryID     int     `json:"category"`
	Size           int     `json:"size"`
	Description    string  `json:"description"`
	ImgURL         string  `json:"img_url"`
	Accepted       bool    `json:"accepted, omitempty"`
	IngredientsIDs []int   `json:"ingredients"`
	Price          float64 `json:"price, omitempty"`
}

// Ingredient type for pizza
type Ingredient struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImgURL      string  `json:"img_url"`
	Price       float64 `json:"price"`
	Created     string  `json:"created_date"`
	Author      User    `json:"author, omitempty"`
}

// IngredientSimple type for all pizzas requests
type IngredientSimple struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	ImgURL string  `json:"img_url"`
	Price  float64 `json:"price"`
}

// UsedIngredient type for ingredients used in custom pizza
type UsedIngredient struct {
	ID           int `json:"id"`
	IngredientID int `json:"ingredient_id"`
	PizzID       int `json:"pizza_id"`
}

// Category is pizzas category - in ex. Vegeterian
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDefault   bool   `json:"is_default"`
}

// CategoryDB is pizzas category as it stored in DB
type CategoryDB struct {
	ID          int
	Name        string
	Description string
	IsDefault   int
}

//Status will be returned as response on signup/login
type Status struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Token      string `json:"token, omitempty"`
}
