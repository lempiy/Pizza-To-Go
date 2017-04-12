package validators

import (
	"regexp"
	"unicode/utf8"

	"html"

	"github.com/lempiy/pizza-app-pq/types"
)

// ValidateEmail validates user email upon signup
func ValidateEmail(email string) bool {
	validEmail := regexp.MustCompile(`^\w+@\w+\.\w{2,4}$`)
	return validEmail.MatchString(email)
}

// ValidatePizza validates user input upon post or put
func ValidatePizza(pizza *types.PizzaPost) (bool, string) {
	if utf8.RuneCountInString(pizza.Name) < 3 {
		return false, "Pizza name is empty or to short"
	}
	if pizza.Size < 30 || pizza.Size > 60 {
		return false, "Incorrect pizza size"
	}
	if len(pizza.IngredientsIDs) < 3 {
		return false, "Pizza ingredients length cannot be lower then 3"
	}
	if pizza.CategoryID == 0 {
		return false, "Category ID cannot be 0"
	}
	if utf8.RuneCountInString(pizza.Description) > 300 {
		return false, "Description length is more then 300 chars"
	}
	return true, ""
}

// EscapeHTMLFromPizza ejects html tags from pizza text fields
// to prevent cross site scripting attacks
func EscapeHTMLFromPizza(pizza *types.PizzaPost) {
	// a value from pointers will be extracted automatically by Go compiler
	pizza.Name = html.EscapeString(pizza.Name)
	pizza.Description = html.EscapeString(pizza.Description)
}
