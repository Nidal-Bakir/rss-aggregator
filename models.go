package main

import (
	"time"

	"github.com/Nidal-Bakir/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

type UserModel struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
}

func toUserModel(dbUser database.User) UserModel {
	return UserModel{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		Token:     dbUser.ApiKey,
	}

}
