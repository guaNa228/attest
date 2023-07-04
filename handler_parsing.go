package main

import (
	"fmt"
	"sync"

	"github.com/guaNa228/attest/logger"
	"github.com/guaNa228/attest/parsing"
)

func (apiCfg *apiConfig) parsingResult() {
	logChan := make(chan string)
	errorChan := make(chan error)

	var errorCounter int

	go logger.Logger(logChan)
	go logger.ErrLogger(errorChan, &errorCounter)

	resultOfParsing := parsing.StartParsing(&logChan, &errorChan, "2023-03-01")

	dbInstances, err := apiCfg.parsingResultToDBInstances(resultOfParsing)
	if err != nil {
		errorChan <- fmt.Errorf("error trying to convert parsed instances to db ones: %s", err)
	}

	if errorCounter > 0 {
		logChan <- "Finishing operation due to insufficient parsing data"
	} else {
		apiCfg.parsedBunkInsert(dbInstances, &logChan, &errorChan, &errorCounter)
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
