package main

import (
	"fmt"
	"net/http"

	db "github.com/guaNa228/attest/internal/database"
)

// Attestation spawn !!!!!!!!!!NEEDED!!!!!!!!!!
// func (apiCfg *apiConfig) handleAttestationSpawn(w http.ResponseWriter, r *http.Request, user db.User) {
// 	type parameters struct {
// 		MonthEnum string `json:"month"`
// 	}
// 	decoder := json.NewDecoder(r.Body)
// 	params := parameters{}
// 	err := decoder.Decode(&params)
// 	if err != nil {
// 		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
// 		return
// 	}

// 	preAttestationData, err := apiCfg.DB.GetPreAttestationData(context.Background())
// 	fmt.Println(len(preAttestationData))
// 	if err != nil {
// 		respondWithError(w, 400, fmt.Sprintf("Database is corrupted: %v", err))
// 	}

// 	attestationData := []*db.Attestation{}

// 	for _, preAttestationItem := range preAttestationData {
// 		attestationData = append(attestationData, &db.Attestation{
// 			ID:                 uuid.New(),
// 			StudentID:          preAttestationItem.StudentID,
// 			SemesterActivityID: preAttestationItem.SemesterActivityID,
// 			Month:              db.MonthEnum(params.MonthEnum),
// 			Result:             sql.NullBool{Valid: false},
// 			Comment:            sql.NullString{Valid: false},
// 		})
// 	}
// 	fmt.Println(len(attestationData))
// 	errorList := itemsBunkCreate(attestationData, "attestation")
// 	if errorList != nil {
// 		for _, err := range errorList {
// 			fmt.Println(err)
// 		}
// 		respondWithJSON(w, 500, "Unable to spawn attestation due to server problems")
// 	}

// 	respondWithJSON(w, 200, struct{}{})
// }

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
	respondWithError(w, 403, "You are not allowed here")
}

// func stubAttestation(n int, index string) []*db.Attestation {
// 	result := []*db.Attestation{}
// 	for i := 0; i < n; i++ {
// 		result = append(result, &Group{
// 			ID:        uuid.New(),
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 			Name:      "test",
// 			Code:      fmt.Sprintf("%v%s", i, index),
// 		})
// 	}
// 	return result

// }
