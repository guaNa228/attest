package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	db "github.com/guaNa228/attest/internal/database"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		Name     string        `json:"name"`
		Login    string        `json:"login"`
		Password string        `json:"password"`
		Role     string        `json:"role"`
		Group_id uuid.NullUUID `json:"group_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	userToCreate, err := apiCfg.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Login:     params.Login,
		Password:  params.Password,
		Role:      params.Role,
		GroupID:   params.Group_id,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	respondWithJSON(w, 201, userToCreate)
}

func (apiCfg *apiConfig) handlerCreateTeacher(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		Name      string `json:"name"`
		TeacherID *int32 `json:"teacher_id,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	//Добавить сравнение по teacher_id
	generatedCredentials, err := apiCfg.credentialsByName(params.Name)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	teacherId := sql.NullInt32{}
	if params.TeacherID != nil {
		teacherId = sql.NullInt32{Valid: true, Int32: *params.TeacherID}
	} else {
		teacherId = sql.NullInt32{Valid: false, Int32: 0}
	}

	fmt.Println(teacherId)

	userToCreate, err := apiCfg.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Login:     generatedCredentials.login,
		Password:  generatedCredentials.password,
		Role:      "teacher",
		GroupID:   uuid.NullUUID{},
		TeacherID: teacherId,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create teacher: %v", err))
		return
	}

	respondWithJSON(w, 200, userToCreate)
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user db.User) {
	grps := stubGroups(5000, "a")
	itemsBunkCreate(grps, "groups")
	respondWithJSON(w, 200, struct{}{})
}

func stubGroups(n int, index string) []*Group {
	result := []*Group{}
	for i := 0; i < n; i++ {
		result = append(result, &Group{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "test",
			Code:      fmt.Sprintf("%v%s", i, index),
		})
	}
	return result
}
