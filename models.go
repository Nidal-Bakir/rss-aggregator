package main

import (
	"database/sql"
	"time"

	"github.com/Nidal-Bakir/rss-aggregator/internal/database"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
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

type FeedModel struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func toFeedModel(feed database.Feed) FeedModel {
	return FeedModel{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
	}
}

type PublicFeedModel struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
}

type FeedFollowsModel struct {
	ID        uuid.UUID `json:"id"`
	FeedId    uuid.UUID `json:"feed_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
}

type PostModel struct {
	ID          uuid.UUID   `json:"id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Title       string      `json:"title"`
	Url         string      `json:"url"`
	PubDate     null.String `json:"pub_date"`
	Description string      `json:"description"`
	FeedID      uuid.UUID   `json:"feed_id"`
}

func toPostModel(post database.Post) PostModel {

	PubDate := null.String{
		NullString: sql.NullString{
			String: post.PubDate.Time.Format(time.RFC3339),
			Valid:  post.PubDate.Valid,
		},
	}

	return PostModel{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Url:         post.Url,
		PubDate:     PubDate,
		Description: post.Description,
		FeedID:      post.FeedID,
	}
}
