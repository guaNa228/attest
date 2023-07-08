package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/jwtauth"
	db "github.com/guaNa228/attest/internal/database"
	"github.com/guaNa228/attest/parsing"
)

type JWT struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Login    string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	var user db.User

	if parsing.IsValidEmail(params.Login) {
		user, err = apiCfg.DB.GetFullUserByEmail(r.Context(), sql.NullString{String: params.Login, Valid: true})
	} else {
		user, err = apiCfg.DB.GetUserByCredentials(r.Context(), db.GetUserByCredentialsParams{
			Login:    params.Login,
			Password: params.Password,
		})
	}

	if err != nil {
		respondWithError(w, 401, "Wrong login or password")
		return
	}

	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": user.ID, "role": user.Role})

	respondWithJSON(w, 201, JWT{
		Token: tokenString,
		Role:  user.Role,
	})

}
