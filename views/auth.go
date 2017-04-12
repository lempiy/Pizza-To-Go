package views

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lempiy/pizza-app-pq/models"
	"github.com/lempiy/pizza-app-pq/sessions"
	"github.com/lempiy/pizza-app-pq/types"
	"github.com/lempiy/pizza-app-pq/utils/utils"
	"github.com/lempiy/pizza-app-pq/utils/validators"
)

//RequiresLogin - middleware to check login state in sessionStorage
func RequiresLogin(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sessions.IsLoggedIn(r) {
			AnswerNonAuthorized(w)
			return
		}
		token := r.Header.Get("Authorization")
		log.Println("validating :" + token)
		if isValid, _ := ValidateToken(token); !isValid {
			AnswerNonAuthorized(w)
			return
		}
		handler(w, r)
	}
}

//Logout Implements the logout functionality. WIll delete the session information from the cookie store
func Logout(w http.ResponseWriter, r *http.Request) {
	var status types.Status
	var message = "Logout successful"
	var token string
	htStatus := http.StatusOK

	session, err := sessions.Store.Get(r, "session")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		htStatus = http.StatusInternalServerError
		message = "Could not logout"
	} else {
		if session.Values["loggedin"] != "false" {
			session.Values["loggedin"] = "false"
			session.Save(r, w)
		}
	}
	w.WriteHeader(htStatus)
	err = json.NewEncoder(w).Encode(status)

	status = types.Status{
		StatusCode: htStatus,
		Message:    message,
		Token:      token}
	if err != nil {
		panic(err)
	}
}

type loginStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//Login implements the login functionality, will add a cookie to the cookie store for managing authentication
func Login(w http.ResponseWriter, r *http.Request) {
	var userData loginStruct

	errParse := json.NewDecoder(r.Body).Decode(&userData)
	defer r.Body.Close()
	if errParse != nil {
		log.Println("Error during parsing userdata\n", errParse)
	}

	var status types.Status
	var message string
	var htStatus = http.StatusOK
	var token string

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	session, err := sessions.Store.Get(r, "session")

	if r.Method == "GET" {
		value := (session.Values["loggedin"] == "true")
		w.WriteHeader(htStatus)
		type loggedIn struct {
			Loggedin bool `json:"loggedin"`
		}
		logged := loggedIn{Loggedin: value}
		json.NewEncoder(w).Encode(logged)
		return
	}

	if errParse != nil {
		log.Println("error while parsing request JSON Body")
		htStatus = http.StatusBadRequest
		message = "Wrong user data"
	} else if err != nil {
		log.Println("error identifying session")
		htStatus = http.StatusInternalServerError
		message = "Could not login"
	} else {
		userData.Password = utils.EncryptPassword(userData.Password)
		//TODO: merge ValidUser and GetUserIDbyName functions
		if (userData.Username != "" && userData.Password != "") && models.ValidUser(userData.Username, userData.Password) {
			id, err := models.GetUserIDbyName(userData.Username)
			session.Values["loggedin"] = "true"
			session.Values["username"] = userData.Username
			session.Values["user_id"] = id
			session.Save(r, w)
			log.Print("user ", userData.Username, " is authenticated")
			token, err = getToken(userData.Username, id)
			if err != nil {
				message = "Failed to getUser token"
				htStatus = http.StatusInternalServerError
			} else {
				message = "Logged in successfully"
			}
		} else {
			htStatus = http.StatusBadRequest
			message = "Invalid user name or password;"
		}
	}

	w.WriteHeader(htStatus)
	status = types.Status{
		StatusCode: htStatus,
		Message:    message,
		Token:      token}
	err = json.NewEncoder(w).Encode(status)

	if err != nil {
		panic(err)
	}
}

type signUpStruct struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//SignUp will enable new users to sign up to our service
func SignUp(w http.ResponseWriter, r *http.Request) {

	var userData signUpStruct
	err := json.NewDecoder(r.Body).Decode(&userData)
	defer r.Body.Close()
	if err != nil {
		log.Println("Error during parsing userdata\n", err)
	}
	var status types.Status
	var message = "Sign up success"
	var statusCode = http.StatusOK
	var token string

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if !checkUniqueUser(userData.Username, userData.Email) {
		log.Printf("User is not unique %v", userData.Username)
		AnswerBadRequest(w, "User with such nickname or email already exists.")
		return
	}

	if !validators.ValidateEmail(userData.Email) {
		log.Printf("User email is not valid %v", userData.Email)
		AnswerBadRequest(w, "User email is not valid.")
		return
	}

	if userData.Username != "" && userData.Password != "" && userData.Email != "" && err == nil {
		userData.Password = utils.EncryptPassword(userData.Password)
		err = models.CreateUser(userData.Username, userData.Password, userData.Email)

		if err != nil {
			log.Println(err)
			statusCode = http.StatusInternalServerError
			message = "Something went wront"
		} else {
			//Login upon signup
			id, err := models.GetUserIDbyName(userData.Username)
			session, _ := sessions.Store.Get(r, "session")
			session.Values["loggedin"] = "true"
			session.Values["username"] = userData.Username
			session.Values["user_id"] = id
			session.Save(r, w)
			log.Print("user ", userData.Username, " is authenticated")
			token, err = getToken(userData.Username, id)
			if err != nil {
				AnswerServerError(w)
				panic(err)
			}
		}
	} else {
		AnswerBadRequest(w, "User data fields cannot be empty.")
		return
	}

	w.WriteHeader(statusCode)
	status = types.Status{
		StatusCode: statusCode,
		Message:    message,
		Token:      token}
	err = json.NewEncoder(w).Encode(status)

	if err != nil {
		panic(err)
	}
}

func checkUniqueUser(username string, email string) bool {
	userID, err := models.GetUserID(username, email)
	if err != nil {
		log.Println(err)
	}
	log.Println(userID)
	if userID > 0 {
		return false
	}
	return true
}
