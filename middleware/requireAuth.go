package middleware

import (
	"context"
	"fmt"
	"github.com/akashc777/csvToPdf/helpers"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"time"
)

func VerifyAuth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tsc, err := r.Cookie("Authorization")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenString := tsc.Value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			hmacSampleSecret := []byte(os.Getenv("JWT_SECRET"))
			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return hmacSampleSecret, nil
		})
		if err != nil {
			log.Fatal(err)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {

			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			userEmail := claims["email"].(string)
			ctx := context.WithValue(r.Context(), helpers.UserEmailContextKey, userEmail)
			r = r.WithContext(ctx)

		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
