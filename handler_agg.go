package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Xzyolanda/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid time_between_reqs %q: %w", cmd.Args[0], err)
	}
	if timeBetweenRequests <= 0 {
		return fmt.Errorf("time_between_reqs must be greater than zero")
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			log.Printf("scrape failed: %v", err)
		}
	}
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	feed, err := s.db.GetNextFeedToFetch(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("couldn't get next feed to fetch: %w", err)
	}

	fmt.Printf("Fetching feed %q: %s\n", feed.Name, feed.Url)
	if err := s.db.MarkFeedFetched(ctx, feed.ID); err != nil {
		return fmt.Errorf("couldn't mark feed %q fetched: %w", feed.Name, err)
	}

	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed %q: %w", feed.Name, err)
	}

	for _, item := range rssFeed.Channel.Item {
		postURL := strings.TrimSpace(item.Link)
		if postURL == "" {
			log.Printf("skipping post with empty URL in feed %q", feed.Name)
			continue
		}

		publishedAt, err := parsePublishedAt(item.PubDate)
		if err != nil {
			log.Printf("couldn't parse published date %q for post %q: %v", item.PubDate, item.Title, err)
		}

		now := time.Now().UTC()
		err = s.db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       item.Title,
			Url:         postURL,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			log.Printf("couldn't save post %q from feed %q: %v", item.Title, feed.Name, err)
		}
	}

	return nil
}

func parsePublishedAt(value string) (sql.NullTime, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return sql.NullTime{}, nil
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC850,
		time.ANSIC,
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 MST",
	}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return sql.NullTime{Time: parsed.UTC(), Valid: true}, nil
		}
	}

	utcLayouts := []string{
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"Mon, 2 Jan 2006 15:04:05",
	}
	for _, layout := range utcLayouts {
		parsed, err := time.ParseInLocation(layout, value, time.UTC)
		if err == nil {
			return sql.NullTime{Time: parsed, Valid: true}, nil
		}
	}

	return sql.NullTime{}, fmt.Errorf("unsupported date format")
}
