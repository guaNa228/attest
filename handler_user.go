package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	db "github.com/guaNa228/attest/internal/database"
	"github.com/guaNa228/attest/translit"
	"github.com/sethvargo/go-password/password"
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

type Credentials struct {
	login    string
	password string
}

func (apiCfg *apiConfig) credentialsByName(fullName string) (Credentials, error) {
	splittedName := strings.Split(translit.ToLatin(strings.ToLower(fullName)), " ")
	log.Println(fullName)
	log.Println(len(fullName))
	log.Println(splittedName)
	var login string
	for index, name := range splittedName {
		if index == 0 {
			login += name + "."
			continue
		}
		login += string(name[0])
	}

	password, err := password.Generate(7, 2, 0, false, true)
	if err != nil {
		return Credentials{}, fmt.Errorf("failed to generate credentials: %s", err)
	}

	isLoginDuplicated, err := apiCfg.DB.IfLoginDuplicates(context.Background(), login)
	if err != nil {
		return Credentials{}, fmt.Errorf("failed to generate credentials: %s", err)
	}

	if isLoginDuplicated {
		numberOfDuplicates, err := apiCfg.DB.NumberOfDuplicatedUsers(context.Background(), login)
		if err != nil {
			return Credentials{}, fmt.Errorf("failed to generate credentials: %s", err)
		}
		login = fmt.Sprintf("%s%v", login, numberOfDuplicates+1)
	}

	return Credentials{
		login:    login,
		password: password,
	}, nil
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

func (apiCfg *apiConfig) handlerCreateStudent(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		Name    string        `json:"name"`
		GroupID uuid.NullUUID `json:"group_id,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	generatedCredentials, err := apiCfg.credentialsByName(params.Name)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	userToCreate, err := apiCfg.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Login:     generatedCredentials.login,
		Password:  generatedCredentials.password,
		Role:      "student",
		GroupID:   params.GroupID,
		TeacherID: sql.NullInt32{},
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create teacher: %v", err))
		return
	}

	respondWithJSON(w, 200, userToCreate)
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user db.User) {
	//Old stubGroups
	// grps := stubGroups(5000, "a")
	// itemsBunkCreate(grps, "groups")
	// respondWithJSON(w, 200, struct{}{})
}

//Old
// func stubGroups(n int, index string) []*db.Group {
// 	result := []*db.Group{}
// 	for i := 0; i < n; i++ {
// 		result = append(result, &db.Group{
// 			ID:        uuid.New(),
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 			S:      "test",
// 			Subcode:      fmt.Sprintf("%v%s", i, index),
// 		})
// 	}
// 	return result
// }
