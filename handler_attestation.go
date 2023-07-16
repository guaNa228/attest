package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	db "github.com/guaNa228/attest/internal/database"
	"github.com/guaNa228/attest/logger"
	"github.com/xuri/excelize/v2"
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
			Result:   sql.NullInt32{Valid: false},
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
		attestationData, err := apiCfg.DB.GetTeachersAttestationData(r.Context(), user.ID)

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

func (apiCfg *apiConfig) handleGetAttestationByWorkload(w http.ResponseWriter, r *http.Request, user db.User) {
	workloadId := chi.URLParam(r, "id")
	workloadUUID, err := uuid.Parse(workloadId)
	if err != nil {
		respondWithError(w, 400, "Broken id")
		return
	}

	attestationData, err := apiCfg.DB.GetWorkloadAttestationData(r.Context(), workloadUUID)
	if err != nil {
		respondWithError(w, 400, "Unknown workload id")
		return
	}

	respondWithJSON(w, 200, attestationData)
}

func (apiCfg *apiConfig) handleAttestationPost(w http.ResponseWriter, r *http.Request, user db.User) {
	type attestationUnit struct {
		Id      uuid.UUID `json:"id"`
		Result  string    `json:"result"`
		Comment string    `json:"comment"`
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
		result := sql.NullInt32{}
		if attestationToUpdate.Result != "null" {
			attestationNumber, err := strconv.Atoi(attestationToUpdate.Result)
			if err != nil {
				respondWithError(w, 400, fmt.Sprintf("Invalid attestation result: %s", attestationToUpdate.Result))
				return
			}
			result = sql.NullInt32{Valid: true, Int32: int32(attestationNumber)}
		} else if attestationToUpdate.Result == "null" {
			result = sql.NullInt32{Valid: false}
		} else {
			respondWithError(w, 400, "Error parsing JSON: wrong result value")
			return
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

func (apiCfg *apiConfig) handleGetUnderachievers(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		StreamID uuid.UUID    `json:"stream"`
		Score    int          `json:"score"`
		Classes  int          `json:"classes"`
		Month    db.MonthEnum `json:"month"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	underachieversData, err := apiCfg.DB.GetUnderachieversData(r.Context(), db.GetUnderachieversDataParams{
		Stream:   params.StreamID,
		Result:   sql.NullInt32{Valid: true, Int32: int32(params.Score)},
		MinScore: params.Classes,
		Month:    params.Month,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Unable to get underachievers data: %v", err))
		return
	}

	if len(underachieversData) == 0 {
		respondWithError(w, 400, "Отсутствуют студенты, соответствующие требованиям")
		return
	}

	sheetName, _ := apiCfg.DB.GetStreamByID(r.Context(), params.StreamID)

	f := excelize.NewFile()

	err = f.SetSheetName("Sheet1", sheetName)
	if err != nil {
		respondWithError(w, 500, "Error deleting default sheet")
		return
	}

	f.SetActiveSheet(0)

	errColA := f.SetColWidth(sheetName, "A", "A", 15)
	errColB := f.SetColWidth(sheetName, "B", "B", 40)
	errColC := f.SetColWidth(sheetName, "C", "C", 50)
	errColD := f.SetColWidth(sheetName, "D", "D", 5)

	if errColA != nil || errColB != nil || errColC != nil || errColD != nil {
		respondWithError(w, 500, "Error customizing file structure")
		return
	}

	currentGroupCode := ""
	currentStudent := ""
	rowIndex := 1
	for _, row := range underachieversData {
		if currentGroupCode != row.GroupCode {
			currentGroupCode = row.GroupCode
			err = f.SetSheetRow(sheetName, fmt.Sprintf("A%v", rowIndex), &[]interface{}{currentGroupCode})
			if err != nil {
				respondWithError(w, 500, fmt.Sprintf("Error filling sheet with group data at row %v: %s", rowIndex, err.Error()))
				return
			}
			rowIndex++
		}

		if currentStudent != row.Student {
			currentStudent = row.Student
			err = f.SetSheetRow(sheetName, fmt.Sprintf("B%v", rowIndex), &[]interface{}{currentStudent})
			if err != nil {
				respondWithError(w, 500, fmt.Sprintf("Error filling sheet with student data at row %v: %s", rowIndex, err.Error()))
				return
			}
			rowIndex++
		}

		err = f.SetSheetRow(sheetName, fmt.Sprintf("C%v", rowIndex), &[]interface{}{row.Class, row.Res.Int32})

		if err != nil {
			respondWithError(w, 500, fmt.Sprintf("Error filling sheet with students data at row %v: %s", rowIndex, err.Error()))
			return
		}

		rowIndex++
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=Неаттестованные_%s_%v.xlsx", sheetName, params.Month))

	err = f.Write(w)
	if err != nil {
		respondWithError(w, 500, "Error generating file")
	}
}
