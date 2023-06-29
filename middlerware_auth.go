package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/guaNa228/attest/internal/auth"
	db "github.com/guaNa228/attest/internal/database"
	"golang.org/x/exp/slices"
)

type authedHandler func(http.ResponseWriter, *http.Request, db.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler, roles []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtToken, err := auth.GetJWT(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error %v", err))
			return
		}

		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			respondWithError(w, 403, fmt.Sprintf("Auth error %v", err))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			respondWithError(w, 403, fmt.Sprintf("Unathorized %v", err))
			return
		}

		user_uuid, err := uuid.Parse(claims["user_id"].(string))

		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Corrupted token data: %v", err))
		}

		user, err := cfg.DB.GetUserById(r.Context(), user_uuid)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Token does not contains valid user info %v", err))
			return
		}

		if user.Role != claims["role"].(string) {
			respondWithError(w, 400, "Corrupted token data")
			return
		}

		if user.Role != "admin" {
			if !slices.Contains(roles, user.Role) {
				respondWithError(w, 403, "You are not allowed here!")
				return
			}
		}

		handler(w, r, user)
	}
}
