package main

import (
	"log"
	"net/http"
	"os"

	"github.com/lempiy/pizza-app-pq/views"
)

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8001"
	}
	views.PopulateTemplates()

	http.HandleFunc("/login/", views.Login)
	http.HandleFunc("/signup/", views.SignUp)
	http.HandleFunc("/logout/", views.Logout)
	http.HandleFunc("/api/categories/", views.RequiresLogin(views.GetCategories))
	http.HandleFunc("/api/ingredients/", views.RequiresLogin(views.GetIngredients))
	http.HandleFunc("/api/save-pizza/", views.RequiresLogin(views.PostPizza))
	http.HandleFunc("/api/accept-pizza/", views.RequiresLogin(views.AcceptPizza))
	http.HandleFunc("/api/delete-pizza/", views.RequiresLogin(views.DeletePizza))
	http.HandleFunc("/api/get-pizzas/", views.GetPizzas)
	http.HandleFunc("/api/get-pizza-by-id/", views.RequiresLogin(views.GetPizzaByID))
	http.HandleFunc("/api/update-pizza/", views.RequiresLogin(views.UpdatePizza))
	http.HandleFunc("/create", views.SendIndex)
	http.HandleFunc("/edit/", views.SendIndex)
	http.Handle("/upload/", http.StripPrefix("/upload/", http.FileServer(http.Dir("./upload"))))
	http.Handle("/", http.FileServer(http.Dir("./dist")))

	log.Print("Running server on " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
