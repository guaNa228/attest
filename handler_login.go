package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/jwtauth"
	db "github.com/guaNa228/attest/internal/database"
)

type JWT struct {
	Token string `json:"token"`
}

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.GetUser(r.Context(), db.GetUserParams{
		Login:    params.Login,
		Password: params.Password,
	})

	if err != nil {
		respondWithError(w, 401, "Wrong login or password")
	}

	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": user.ID, "role": user.Role})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
	respondWithJSON(w, 201, JWT{
		Token: tokenString,
	})

}
