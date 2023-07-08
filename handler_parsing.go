package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	db "github.com/guaNa228/attest/internal/database"
	"github.com/guaNa228/attest/logger"
	"github.com/guaNa228/attest/parsing"
)

func (apiCfg *apiConfig) handleParsing(w http.ResponseWriter, r *http.Request, user db.User) {

	decoder := json.NewDecoder(r.Body)
	params := parsing.Parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	var errorCounter int

	apiCfg.parsingResult(params, &errorCounter)

	if errorCounter > 0 {
		respondWithError(w, 400, "Something went wrong, see the error log")
	} else {
		respondWithJSON(w, 201, struct{}{})
	}

	GlobalWsWg.Done()
}

func (apiCfg *apiConfig) parsingResult(params parsing.Parameters, errCounter *int) {
	logChan := make(chan string)
	errorChan := make(chan error)

	go logger.Logger(logChan, GlobalWsConn, true)
	go logger.ErrLogger(errorChan, errCounter, GlobalWsConn, true)

	resultOfParsing := parsing.StartParsing(&logChan, &errorChan, "2023-03-01", params)

	dbInstances, err := apiCfg.parsingResultToDBInstances(resultOfParsing)
	if err != nil {
		errorChan <- fmt.Errorf("error trying to convert parsed instances to db ones: %s", err)
	}

	if *errCounter > 0 {
		logChan <- "Finishing operation due to insufficient parsing data"
		return
	} else {
		apiCfg.parsedBunkInsert(dbInstances, &logChan, &errorChan, errCounter)
	}

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
}
