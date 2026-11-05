package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Xzyolanda/blog_aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s [limit]", cmd.Name)
	}

	limit := int32(2)
	if len(cmd.Args) == 1 {
		parsedLimit, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil || parsedLimit <= 0 {
			return fmt.Errorf("limit must be a positive integer")
		}
		limit = int32(parsedLimit)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for current user: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\nURL: %s\nFeed: %s\n", post.Title, post.Url, post.FeedName)
		if post.PublishedAt.Valid {
			fmt.Printf("Published: %s\n", post.PublishedAt.Time.Format(time.RFC3339))
		}
		if post.Description.Valid {
			fmt.Printf("Description: %s\n", post.Description.String)
		}
		fmt.Println()
	}

	return nil
}
