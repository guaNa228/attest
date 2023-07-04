package main

import (
	"context"
	"errors"

	"github.com/google/uuid"
	db "github.com/guaNa228/attest/internal/database"
)

func (apiCfg *apiConfig) createPrograms() error {
	programsFullTime := map[string]int16{
		"Бакалавриат":  4,
		"Магистратура": 2,
		"Специалитет":  5,
	}

	programsNumber, err := apiCfg.DB.GetProgramsNumber(context.Background())
	if err != nil {
		return err
	}

	if programsNumber == int64(len(programsFullTime)) {
		return errors.New("programs are already initialized")
	}

	for name, maxCourses := range programsFullTime {
		_, err := apiCfg.DB.CreateProgram(context.Background(), db.CreateProgramParams{
			ID:         uuid.New(),
			Name:       name,
			MaxCourses: maxCourses,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
