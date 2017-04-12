package views

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lempiy/pizza-app-pq/utils/utils"
)

//TokenClaims used for JWT
type TokenClaims struct {
	Username string `json:"username"`
	ID       int    `json:"user_id"`
	jwt.StandardClaims
}

var mySigningKey = []byte(utils.Keys.TokenKey)

//GetTokenHandler will get a token for the username
func getToken(username string, id int) (string, error) {

	claims := TokenClaims{
		username,
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 5).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	/* Sign the token with secret */
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//ValidateToken will validate the token
func ValidateToken(incomingToken string) (bool, string) {
	token, err := jwt.ParseWithClaims(incomingToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})

	if err != nil {
		return false, ""
	}

	claims := token.Claims.(*TokenClaims)
	return token.Valid, claims.Username
}
