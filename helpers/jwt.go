package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func DecodeTokenJwt(tokenStr string) (jwt.MapClaims, bool) {

	hmacSecretString := os.Getenv("JWT_SECRET")
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		fmt.Println("Invalid Jwt")
		return nil, false
	}
}

func GenerateTokenAuth(Name, Email string, UserId, UserTypes int64) string {
	token := jwt.New(jwt.SigningMethodHS256)
	value := token.Claims.(jwt.MapClaims)

	value["Name"] = Name
	value["Email"] = Email
	value["UserId"] = UserId
	value["UserTypes"] = UserTypes
	value["Expired"] = time.Now().Add(time.Hour * 1).Format("2006-01-02 15:04:05")

	jwtKey := os.Getenv("JWT_SECRET")

	tokenString, _ := token.SignedString([]byte(jwtKey))

	return tokenString
}