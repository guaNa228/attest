package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	db "github.com/guaNa228/attest/internal/database"
)

func (apiCfg *apiConfig) handlerCreateClass(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	classToCreate, err := apiCfg.DB.CreateClass(r.Context(), db.CreateClassParams{
		ID:   uuid.New(),
		Name: params.Name,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create class: %v", err))
		return
	}

	respondWithJSON(w, 201, classToCreate)
}

func (apiCfg *apiConfig) handlerDeleteClass(w http.ResponseWriter, r *http.Request, user db.User) {
	const instance = "class"
	const paramToSearch = "classToDeleteID"
	classToDelete := chi.URLParam(r, paramToSearch)
	if classToDelete == "" {
		respondWithError(w, 400, fmt.Sprintf("Wrong request address. Should be {%v}/{%v}, not {%v}?{%v}Id={%v}",
			instance,
			paramToSearch,
			instance,
			instance,
			paramToSearch))
	}

	classToDeleteID, err := uuid.Parse(classToDelete)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Corrupted %v id: %v", instance, err))
		return
	}

	err = apiCfg.DB.DeleteClassByID(r.Context(), classToDeleteID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete %v by ID: %v", instance, err))
		return
	}

	respondWithJSON(w, 200, struct{}{})
}
