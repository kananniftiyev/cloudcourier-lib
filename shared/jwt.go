package shared

import (
	"context"
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

const SECRET_KEY = "secret"

type CustomClaims struct {
	jwt.StandardClaims
	UserID   uint
	Username string

	// Add other custom claims as needed
}

// JWTTokenVerifyMiddleware is a middleware function for verifying JWT tokens
func JWTTokenVerifyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the JWT token from the cookie
		cookie, err := r.Cookie("jwt")
		if err != nil {
			RespondWithError(w, err, http.StatusForbidden)
			return
		}

		// Verify the JWT token
		token, err := jwt.ParseWithClaims(cookie.Value, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		})

		if err != nil {
			RespondWithError(w, err, http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			RespondWithError(w, err, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*CustomClaims)

		if !ok {
			RespondWithError(w, errors.New("failed to get token claims"), http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
