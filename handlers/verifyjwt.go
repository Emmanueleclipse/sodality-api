package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"sodality/models"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var JWT_SECRET = []byte(DotEnvVariable("JWT_SECRET"))

type Claims struct {
	Username string `json:"username"`
	Dash     string `json:"dash"`
	// Email    string `json:"email"`
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// IsAuthorized -> verify jwt header
func IsAuthorized(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")

		if len(authHeader) != 2 {
			AuthorizationResponse("Malformed JWT token", w)
		} else {
			jwtToken := authHeader[1]
			token, _ := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return JWT_SECRET, nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				ctx := context.WithValue(r.Context(), "props", claims)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				AuthorizationResponse("Unauthorized", w)
			}
		}
	})
}

// GenerateJWT -> generate jwt
func GenerateJWT(user models.User) (string, error) {
	claims := &Claims{
		Username: user.Username,
		Dash:     user.Dash,
		// Email:    user.Email,
		UserID: user.ID.Hex(),
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JWT_SECRET)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
