package views

import (
	"encoding/json"
	"log"
	"net/http"

	"fmt"

	"strconv"

	"github.com/lempiy/pizza-app-pq/models"
	"github.com/lempiy/pizza-app-pq/sessions"
	"github.com/lempiy/pizza-app-pq/types"
	"github.com/lempiy/pizza-app-pq/utils/utils"
	"github.com/lempiy/pizza-app-pq/utils/validators"
)

//GetPizzas used to send non-deleted pizzas to user
func GetPizzas(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		AnswerBadRequest(w, "Wrong request method only GET allowed.")
		return
	}

	statusCode := http.StatusOK
	pizzas, err := models.GetAllPizzas()

	if err != nil {
		log.Printf("Error during getting pizzas for DB.")
		statusCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	err = json.NewEncoder(w).Encode(pizzas)

	if err != nil {
		panic(err)
	}
}

//GetCategories used to send categories to user
func GetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		AnswerBadRequest(w, "Wrong request method only GET allowed.")
		return
	}

	statusCode := http.StatusOK
	categories, err := models.GetAllCategories()

	if err != nil {
		log.Printf("Error during getting categories for DB.")
		statusCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(categories)

	if err != nil {
		panic(err)
	}
}

//GetIngredients - handler for getting all existing ingredients
func GetIngredients(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		AnswerBadRequest(w, "Wrong request method only GET allowed.")
		return
	}

	statusCode := http.StatusOK
	ingredients, err := models.GetAllIngredients()

	if err != nil {
		log.Printf("Error during getting categories for DB.")
		statusCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(ingredients)

	if err != nil {
		panic(err)
	}
}

// PizzaRequest is a struct to extract user post form
type PizzaRequest struct {
	Name           string `json:"name"`
	CategoryID     int    `json:"category"`
	Size           int    `json:"size"`
	Description    string `json:"description"`
	EncodedImage   string `json:"encodedImage"`
	IngredientsIDs []int  `json:"ingredients"`
}

// PostAnswer - type for answer
type PostAnswer struct {
	Success bool `json:"success"`
}

// PostPizza main handler to push new pizza orders
func PostPizza(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		AnswerBadRequest(w, "Wrong request method only POST allowed.")
		return
	}
	var pizzaAnswer PostAnswer
	var pizzaReq PizzaRequest

	err := json.NewDecoder(r.Body).Decode(&pizzaReq)
	log.Println("Save image")
	imgName, err := utils.SaveEncodedImage(pizzaReq.EncodedImage)
	log.Println("Saved image " + imgName)
	if err != nil {
		log.Println("Saved image error " + imgName)
		AnswerServerError(w)
		panic(err)
	}

	pizzaPost := types.PizzaPost{
		Name:           pizzaReq.Name,
		AuthorID:       sessions.GetCurrentUserID(r),
		CategoryID:     pizzaReq.CategoryID,
		Size:           pizzaReq.Size,
		Description:    pizzaReq.Description,
		ImgURL:         "upload/" + imgName,
		IngredientsIDs: pizzaReq.IngredientsIDs}

	if isOk, errMessage := validators.ValidatePizza(&pizzaPost); !isOk {
		AnswerBadRequest(w, errMessage)
		return
	}

	validators.EscapeHTMLFromPizza(&pizzaPost)

	pizzaPrice, err := models.CalculatePizzaPrice(pizzaPost.IngredientsIDs, pizzaPost.Size)

	if err != nil {
		AnswerServerError(w)
		panic(err)
	}

	err = models.PostPizza(pizzaPost, pizzaPrice)
	pizzaAnswer.Success = true
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(pizzaAnswer)
	if err != nil {
		AnswerServerError(w)
		panic(err)
	}
}

type updatePizzaReq struct {
	ID int `json:"id"`
}
type answer struct {
	Message string `json:"message"`
}

// DeletePizza - marks pizza as deleted
func DeletePizza(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		AnswerBadRequest(w, "Wrong request method only PUT allowed.")
		return
	}

	var pizzaReq updatePizzaReq
	var userID int
	var anwr answer
	err := json.NewDecoder(r.Body).Decode(&pizzaReq)
	if err != nil {
		AnswerServerError(w)
		panic(err)
	}
	if pizzaReq.ID != 0 {
		userID = sessions.GetCurrentUserID(r)
		log.Printf("Put pizzaID - %d | from user - %d\n", pizzaReq.ID, userID)
		count, err := models.DeletePizza(pizzaReq.ID, userID)
		if count == 0 {
			AnswerBadRequest(w, "Wrong pizza ID or non-auth access")
			return
		}
		if err != nil {
			AnswerServerError(w)
			panic(err)
		}
		anwr.Message = fmt.Sprintf("Successfully deleted pizza with #%d", pizzaReq.ID)
		err = json.NewEncoder(w).Encode(anwr)
	} else {
		AnswerBadRequest(w, "Pizza ID cannot be 0")
		return
	}
}

// AcceptPizza - marks pizza as accepted
func AcceptPizza(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		AnswerBadRequest(w, "Wrong request method only PUT allowed.")
		return
	}

	var pizzaReq updatePizzaReq
	var userID int
	var anwr answer
	err := json.NewDecoder(r.Body).Decode(&pizzaReq)
	if err != nil {
		AnswerServerError(w)
		panic(err)
	}
	if pizzaReq.ID != 0 {
		userID = sessions.GetCurrentUserID(r)
		log.Printf("Put pizzaID - %d | from user - %d\n", pizzaReq.ID, userID)
		count, err := models.AcceptPizza(pizzaReq.ID, userID)
		if count == 0 {
			AnswerBadRequest(w, "Wrong pizza ID or non-auth access")
			return
		}
		if err != nil {
			AnswerServerError(w)
			panic(err)
		}
		anwr.Message = fmt.Sprintf("Successfully accepted pizza with #%d", pizzaReq.ID)
		err = json.NewEncoder(w).Encode(anwr)
	} else {
		AnswerBadRequest(w, "Pizza ID cannot be 0")
		return
	}
}

type pizzaByIDAnswer struct {
	IsMyPizza bool        `json:"is_my_pizza"`
	Pizza     types.Pizza `json:"pizza"`
}

//GetPizzaByID used to get pizza by ID and send user message
//about whether or not this pizza is belongs to him.
func GetPizzaByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		AnswerBadRequest(w, "Wrong request method only GET allowed.")
		return
	}
	var pizzaByID pizzaByIDAnswer
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		AnswerBadRequest(w, "Wrong id param, expected to be a positive integer and not 0")
		return
	}
	statusCode := http.StatusOK
	log.Println("PIZZA ID:", id)
	pizza, err := models.GetPizzaByID(id)
	requesterID := sessions.GetCurrentUserID(r)
	log.Println("REQUESTER ID:", id)
	pizzaByID.IsMyPizza = pizza.ID != 0 && pizza.Author.ID == requesterID
	pizzaByID.Pizza = pizza

	if err != nil {
		log.Printf("Error during getting pizzas for DB.")
		statusCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(pizzaByID)

	if err != nil {
		panic(err)
	}
}

// UpdatePizza updates currently edited pizza by ID
func UpdatePizza(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		AnswerBadRequest(w, "Wrong request method only PUT allowed.")
		return
	}
	var pizzaAnswer PostAnswer
	var pizzaReq PizzaRequest

	//check if pizza is users
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		AnswerBadRequest(w, "Wrong id param, expected to be a positive integer and not 0")
		return
	}

	pizza, err := models.GetPizzaByID(id)
	requesterID := sessions.GetCurrentUserID(r)
	log.Printf("My pizza author - %d | User id - %d\n", pizza.Author.ID, requesterID)
	if requesterID != pizza.Author.ID {
		AnswerBadRequest(w, "Thats not your pizza, you cannot edit it.")
		return
	}

	err = json.NewDecoder(r.Body).Decode(&pizzaReq)

	if err != nil {
		AnswerServerError(w)
		panic(err)
	}

	_, err = models.RemoveUnusedIngredients(id)

	if err != nil {
		AnswerServerError(w)
		panic(err)
	}

	err = utils.RemoveUnusedImg(pizza.ImgURL)
	if err != nil {
		log.Println(err)
	}

	imgName, err := utils.SaveEncodedImage(pizzaReq.EncodedImage)

	if err != nil {
		AnswerServerError(w)
		panic(err)
	}

	pizzaPost := types.PizzaPost{
		Name:           pizzaReq.Name,
		AuthorID:       sessions.GetCurrentUserID(r),
		CategoryID:     pizzaReq.CategoryID,
		Size:           pizzaReq.Size,
		Description:    pizzaReq.Description,
		ImgURL:         "upload/" + imgName,
		IngredientsIDs: pizzaReq.IngredientsIDs}

	if isOk, errMessage := validators.ValidatePizza(&pizzaPost); !isOk {
		AnswerBadRequest(w, errMessage)
		return
	}

	validators.EscapeHTMLFromPizza(&pizzaPost)

	pizzaPrice, err := models.CalculatePizzaPrice(pizzaPost.IngredientsIDs, pizzaPost.Size)

	if err != nil {
		AnswerServerError(w)
		panic(err)
	}

	err = models.UpdatePizza(pizzaPost, pizzaPrice, id, requesterID)
	pizzaAnswer.Success = true
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(pizzaAnswer)
	if err != nil {
		AnswerServerError(w)
		panic(err)
	}
}
