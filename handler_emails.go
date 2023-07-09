package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	db "github.com/guaNa228/attest/internal/database"
	"github.com/guaNa228/attest/logger"
	"github.com/guaNa228/attest/parsing"
)

func (apiCfg *apiConfig) handleEmailParsing(w http.ResponseWriter, r *http.Request, user db.User) {

	logChan := make(chan string)
	errorChan := make(chan error)

	var errorCounter int

	go logger.Logger(logChan, GlobalWsConn, true)
	go logger.ErrLogger(errorChan, &errorCounter, GlobalWsConn, true)

	parsedTeachersEmails := parsing.ParseTeachersMails(apiCfg.DB, &logChan, &errorChan)

	fmt.Println(len(*parsedTeachersEmails))

	updateWg := sync.WaitGroup{}
	updateWg.Add(1)
	itemsBunkUpdate(*parsedTeachersEmails, "users", "email", "id", &updateWg, &errorChan, &errorCounter)

	updateWg.Wait()

	logWG := sync.WaitGroup{}
	logWG.Add(1)
	go func() {
		defer logWG.Done()
		close(logChan)
	}()

	errorWG := sync.WaitGroup{}
	errorWG.Add(1)
	go func() {
		defer errorWG.Done()
		close(errorChan)
	}()

	logWG.Wait()
	errorWG.Wait()

	if errorCounter > 0 {
		respondWithError(w, 400, "Something went wrong, see the error log")
	} else {
		respondWithJSON(w, 201, struct{}{})
	}

	GlobalWsWg.Done()
}

func (apiCfg *apiConfig) handleGetEmails(w http.ResponseWriter, r *http.Request, user db.User) {
	emails, err := apiCfg.DB.GetTeachersEmails(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Unable to get teachers data: %s", err.Error()))
		return
	}

	respondWithJSON(w, 200, emails)
}

func (apiCfg *apiConfig) handleGetStudentsEmails(w http.ResponseWriter, r *http.Request, user db.User) {
	emails, err := apiCfg.DB.GetUsersEmails(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Unable to get students data: %s", err.Error()))
		return
	}
	respondWithJSON(w, 200, emails)
}

func (apiCfg *apiConfig) handlerUpdateMails(w http.ResponseWriter, r *http.Request, user db.User) {

	type parameters struct {
		Mails []parsing.ParsedTeachersEmails `json:"mails"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	failedEmails := []string{}

	for _, emailUnit := range params.Mails {
		err := apiCfg.DB.UpdateEmail(r.Context(), db.UpdateEmailParams{
			ID: emailUnit.Id,
			Email: sql.NullString{
				Valid:  true,
				String: emailUnit.Email,
			},
		})

		if err != nil {
			failedEmails = append(failedEmails, emailUnit.Email)
		}

	}

	if len(failedEmails) == 0 {
		respondWithJSON(w, 200, "All emails are added")
	} else {
		respondWithError(w, 400, fmt.Sprintf("Following emails are not added: %s", strings.Join(failedEmails, ",")))
	}
}
