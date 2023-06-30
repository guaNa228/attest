package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	db "github.com/guaNa228/attest/internal/database"
)

func (apiCfg *apiConfig) handlerCreateGroup(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		Name string `json:"name"`
		Code string `json:"code"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	groupToCreate, err := apiCfg.DB.CreateGroup(r.Context(), db.CreateGroupParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Code:      params.Code,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create group: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseGroupToGroup(groupToCreate))
}

func (apiCfg *apiConfig) handlerDeleteGroup(w http.ResponseWriter, r *http.Request, user db.User) {
	groupToDelete := chi.URLParam(r, "groupToDelete")
	if groupToDelete == "" {
		respondWithError(w, 400, "Wrong request address. Should be group/{groupToDelete}, not group?groupId={groupToDelete}")
	}
	groupToDeleteID, err := uuid.Parse(groupToDelete)
	if err != nil {
		err = apiCfg.DB.DeleteGroupByCode(r.Context(), groupToDelete)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't find group by code: %v", err))
			return
		}
	}

	err = apiCfg.DB.DeleteGroupByID(r.Context(), groupToDeleteID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete group by ID: %v", err))
		return
	}

	respondWithJSON(w, 200, struct{}{})
}
