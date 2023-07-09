package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	db "github.com/guaNa228/attest/internal/database"
	"github.com/guaNa228/attest/logger"
)

func (apiCfg *apiConfig) handleAttestationSpawn(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		MonthEnum string `json:"month"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	preAttestationData, err := apiCfg.DB.GetPreAttestationData(context.Background())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Database is corrupted: %v", err))
	}

	attestationData := []*db.Attestation{}

	for _, preAttestationItem := range preAttestationData {
		attestationData = append(attestationData, &db.Attestation{
			ID:       uuid.New(),
			Student:  preAttestationItem.Student,
			Workload: preAttestationItem.Workload,
			Month:    db.MonthEnum(params.MonthEnum),
			Result:   sql.NullBool{Valid: false},
			Comment:  sql.NullString{Valid: false},
		})
	}

	errorChan := make(chan error)

	var errorCounter int

	go logger.ErrLogger(errorChan, &errorCounter, GlobalWsConn, false)

	outerWg := sync.WaitGroup{}

	outerWg.Add(1)

	itemsBunkCreate(attestationData, "attestation", &outerWg, &errorChan, &errorCounter)

	outerWg.Wait()

	errorWG := sync.WaitGroup{}
	errorWG.Add(1)
	go func() {
		defer errorWG.Done()
		close(errorChan)
	}()

	errorWG.Wait()

	if errorCounter == 0 {
		respondWithJSON(w, 200, struct{}{})
	} else {
		respondWithError(w, 400, fmt.Sprintf("Аттестация за %s уже была начата", params.MonthEnum))
	}

}

func (apiCfg *apiConfig) handleAttestationGet(w http.ResponseWriter, r *http.Request, user db.User) {
	if user.Role == "teacher" {
		attestationData, err := apiCfg.DB.GetAttestationData(r.Context(), user.ID)

		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get attestation data: %v", err))
			return
		}

		respondWithJSON(w, 200, attestationData)
		return
	}
	if user.Role == "student" {
		attestationData, err := apiCfg.DB.GetStudentsAttestationData(r.Context(), user.ID)

		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get attestation data: %v", err))
			return
		}

		respondWithJSON(w, 200, attestationData)
		return
	}
	respondWithError(w, 403, "You are not allowed here")
}

type NullableBool struct {
	Valid bool
	Value bool
}

func (nb *NullableBool) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*nb = NullableBool{Valid: false, Value: false}
		return nil
	} else {
		if string(data) == "true" {
			*nb = NullableBool{Valid: true, Value: true}
			return nil
		} else if string(data) == "false" {
			*nb = NullableBool{Valid: false, Value: true}
			return nil
		} else {
			return errors.New("unknown type passed into result")
		}
	}
}

func (apiCfg *apiConfig) handleAttestationPost(w http.ResponseWriter, r *http.Request, user db.User) {
	type attestationUnit struct {
		Id      uuid.UUID    `json:"id"`
		Result  NullableBool `json:"result"`
		Comment string       `json:"comment"`
	}

	type parameters struct {
		Data []attestationUnit `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	for _, attestationToUpdate := range params.Data {
		result := sql.NullBool{}
		if attestationToUpdate.Result.Valid {
			result = sql.NullBool{Valid: true, Bool: attestationToUpdate.Result.Value}
		}
		comment := sql.NullString{}
		if attestationToUpdate.Comment != "" {
			comment = sql.NullString{Valid: true, String: attestationToUpdate.Comment}
		}
		err := apiCfg.DB.UpdateAttestationRow(r.Context(), db.UpdateAttestationRowParams{
			ID:      attestationToUpdate.Id,
			Result:  result,
			Comment: comment,
		})
		if err != nil {
			respondWithError(w, 400, "Attestation row not found")
			return
		}
	}

	respondWithJSON(w, 200, struct{}{})
}

func (apiCfg *apiConfig) handleAttestationClear(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		Month db.MonthEnum `json:"month"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	err = apiCfg.DB.ClearAttestation(r.Context(), params.Month)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error clearing atestation data for %s", params.Month))
	} else {
		respondWithJSON(w, 200, struct{}{})
	}
}
