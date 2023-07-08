package main

import (
	"fmt"
	"net/http"
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
