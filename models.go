package main

import (
	"time"

	"github.com/ByChanderZap/rss-web-server/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		ApiKey:    dbUser.ApiKey,
	}
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserId    uuid.UUID `json:"user_id"`
}

func databaseFeedToFeed(feedDb database.Feed) Feed {
	return Feed{
		ID:        feedDb.ID,
		CreatedAt: feedDb.CreatedAt,
		UpdatedAt: feedDb.UpdatedAt,
		Name:      feedDb.Name,
		Url:       feedDb.Url,
		UserId:    feedDb.UserID,
	}
}

func dbFeedsToFeeds(dbFeeds []database.Feed) []Feed {
	feeds := []Feed{}
	for _, dbFeed := range dbFeeds {
		feeds = append(feeds, databaseFeedToFeed(dbFeed))
	}
	return feeds
}

type FeedFollows struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    uuid.UUID `json:"user_id"`
	FeedId    uuid.UUID `json:"feed_id"`
}

func dbFFToFF(ffDb database.FeedsFollow) FeedFollows {
	return FeedFollows{
		ID:        ffDb.ID,
		CreatedAt: ffDb.CreatedAt,
		UpdatedAt: ffDb.UpdatedAt,
		UserId:    ffDb.UserID,
		FeedId:    ffDb.FeedID,
	}
}

func dbFeedfsToFeedf(ffsDb []database.FeedsFollow) []FeedFollows {
	ffs := []FeedFollows{}
	for _, dbFeed := range ffsDb {
		ffs = append(ffs, dbFFToFF(dbFeed))
	}
	return ffs
}

type Posts struct {
	Id          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Url         string    `json:"url"`
	FeedId      uuid.UUID `json:"feed_id"`
}

func dbPostToPost(dbPost database.Post) Posts {
	var description *string
	if dbPost.Description.Valid {
		description = &dbPost.Description.String
	}
	return Posts{
		Id:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		Description: description,
		PublishedAt: dbPost.PublishedAt,
		Url:         dbPost.Url,
		FeedId:      dbPost.FeedID,
	}
}

func dbPostsToPosts(dbPosts []database.Post) []Posts {
	posts := []Posts{}

	for _, p := range dbPosts {
		posts = append(posts, dbPostToPost(p))
	}

	return posts
}
