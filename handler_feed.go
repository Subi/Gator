package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/subi/gator/internal/database"
	"html"
	"io"
	"net/http"
	"time"
)

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

func handlerFetchFeed(s *state, cmd Command) error {
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
	fmt.Printf("%s\n", html.UnescapeString(string(data)))
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

func handlerAddFeed(s *state, cmd Command) error {
	if len(cmd.args) < 2 {
		return errors.New("please provide a name and url for feed")
	}
	u, err := s.db.GetUser(context.Background(), s.config.CurrentUsername)
	if err != nil {
		return err
	}
	f, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    u.ID,
	})
	fmt.Print(f)
	return nil
}
