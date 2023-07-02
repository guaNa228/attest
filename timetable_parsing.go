package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/guaNa228/attest/translit"
	"github.com/sethvargo/go-password/password"
)

type Credentials struct {
	login    string
	password string
}

func (apiCfg *apiConfig) credentialsByName(fullName string) (Credentials, error) {
	splittedName := strings.Split(translit.ToLatin(strings.ToLower(fullName)), " ")
	surname, name, fathername := splittedName[0], splittedName[1], splittedName[2]
	login := fmt.Sprintf("%s.%v%v", surname, string(name[0]), string(fathername[0]))

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
