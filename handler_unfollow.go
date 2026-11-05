package main

import (
	"context"
	"fmt"

	"github.com/Xzyolanda/blog_aggregator/internal/database"
)

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	ctx := context.Background()
	feed, err := s.db.GetFeedByURL(ctx, cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't find feed: %w", err)
	}

	err = s.db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't unfollow feed: %w", err)
	}

	fmt.Printf("%s is no longer following %s\n", user.Name, feed.Name)
	return nil
}
