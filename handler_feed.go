package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/google/uuid"
	pq2 "github.com/lib/pq"
	"github.com/subi/gator/internal/database"
	"html"
	"io"
	"net/http"
	"time"
)

const duplicateURLKeyError = "23505"

type Feed struct {
	Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Item        []Item `xml:"item"`
	} `xml:"channel"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func handlerFetchFeed(s *state, cmd Command, url string, feedId uuid.UUID) error {
	c := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	feedContent := Feed{}
	err = xml.Unmarshal(data, &feedContent)
	if err != nil {
		return err
	}
	var pq *pq2.Error
	for _, item := range feedContent.Channel.Item {
		_, err = s.db.CreatPosts(context.Background(), database.CreatPostsParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: escapedString(item.Description),
			PublishedAt: item.PubDate,
			FeedID:      feedId,
		})
		if err != nil {
			errors.As(err, &pq)
			if pq.Code == duplicateURLKeyError {
				continue
			}
			return err
		}
	}
	return nil
}

func handlerListFeeds(s *state, cmd Command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		users, err := s.db.GetUsers(context.Background())
		if err != nil {
			return err
		}
		for _, user := range users {
			if user.ID == feed.UserID {
				fmt.Println(user.Name)
			}
		}
	}
	return nil
}

func handlerAddFeed(s *state, cmd Command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("please provide a name and url for feed")
	}
	name := cmd.args[0]
	url := cmd.args[1]
	f, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    f.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func handlerFollowFeed(s *state, cmd Command, user database.User) error {
	url := cmd.args[0]
	f, err := s.db.FindFeed(context.Background(), url)
	if err != nil {
		return err
	}
	followFeed, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    f.ID,
	})
	fmt.Print(followFeed)
	return nil
}

func handlerFollowing(s *state, cmd Command, user database.User) error {
	following, err := s.db.GetFeedFollowForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, follow := range following {
		fmt.Println(follow.FeedName)
	}
	return nil
}

func handlerUnfollowFeed(s *state, cmd Command, user database.User) error {
	url := cmd.args[0]
	err := s.db.DeleteFeedFollowForUser(context.Background(), database.DeleteFeedFollowForUserParams{
		UserID: user.ID,
		Url:    url,
	})
	if err != nil {
		return err
	}
	return nil
}

func initScrape(s *state, cmd Command) error {
	if len(cmd.args) < 1 {
		return errors.New("please provide a interval")
	}
	duration := cmd.args[0]
	parsedInterval, err := time.ParseDuration(duration)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(parsedInterval)
	for ; ; <-ticker.C {
		fmt.Printf("Collecting feed every %s\n", cmd.args[0])
		err = scrapeFeed(s, cmd)
		if err != nil {
			return err
		}
	}
}

func scrapeFeed(s *state, cmd Command) error {
	f, err := s.db.GetNextFeedFetch(context.Background())
	if err != nil {
		return err
	}
	err = s.db.MarkFeedFetched(context.Background(), f.ID)
	if err != nil {
		return err
	}
	err = handlerFetchFeed(s, cmd, f.Url, f.ID)
	if err != nil {
		return err
	}
	return nil
}

func escapedString(s string) string {
	return html.UnescapeString(s)
}
