package Middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

var signingKey = []byte(os.Getenv("SECRET_CODE"))

// func generateJWT(username string, secretKey []byte) (string, error) {
// 	// Create a new token
// 	token := jwt.New(jwt.SigningMethodHS256)

// 	// Set claims
// 	claims := token.Claims.(jwt.MapClaims)
// 	claims["username"] = username
// 	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiration time (24 hours)

// 	// Sign the token with the secret key
// 	tokenString, err := token.SignedString(secretKey)
// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenString, nil
// }

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Extract JWT token from the Authorization header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return signingKey, nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Check if the token is valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Println("Authorized user:", claims["username"])
			// If the token is valid, call the next handler
			r = r.WithContext(context.WithValue(r.Context(), "username", claims["username"]))
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}
