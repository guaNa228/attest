package main

import (
	"time"

	"github.com/google/uuid"
	db "github.com/guaNa228/attest/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Login     string    `json:"login"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
}

func databaseUserToUser(dbUser db.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		Login:     dbUser.Login,
		Password:  dbUser.Password,
		Role:      dbUser.Role,
	}
}

type Group struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
}

func databaseGroupToGroup(dbGroup db.Group) Group {
	return Group{
		ID:        dbGroup.ID,
		CreatedAt: dbGroup.CreatedAt,
		UpdatedAt: dbGroup.UpdatedAt,
		Name:      dbGroup.Name,
		Code:      dbGroup.Code,
	}
}
