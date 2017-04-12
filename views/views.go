package views

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/lempiy/pizza-app-pq/types"
)

var (
	indexTemplate *template.Template
	templates     *template.Template
)

//PopulateTemplates is used to parse all templates present in
//the templates folder
func PopulateTemplates() {
	var allFiles []string
	templatesDir := "./dist/"
	files, err := ioutil.ReadDir(templatesDir)
	if err != nil {
		fmt.Println("Error reading template dir")
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") {
			allFiles = append(allFiles, templatesDir+filename)
		}
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	templates = template.Must(template.ParseFiles(allFiles...))

	indexTemplate = templates.Lookup("index.html")

}

//SendIndex sends index HTML SPA page to client
func SendIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var context interface{}
		indexTemplate.Execute(w, context)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

//AnswerBadRequest - is used to answer with 400 error
func AnswerBadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	status := types.Status{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		Token:      ""}
	err := json.NewEncoder(w).Encode(status)
	if err != nil {
		log.Println(err)
	}
}

//AnswerNonAuthorized - answers user with no creaditians
func AnswerNonAuthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Wrong Authorization Creaditians."))
}

//AnswerServerError - answers user for server errors
func AnswerServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Server error."))
}
