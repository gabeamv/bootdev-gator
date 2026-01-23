package gatorcommand

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gabeamv/bootdev-gator/gatorfeed"
	"github.com/gabeamv/bootdev-gator/internal/database"
)

// Handler function to handle the login command.
func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("arguments expected: 'username'")
	}
	username := cmd.Args[0]
	user, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("error getting the user '%v': %w", username, err)
	}
	err = s.S.SetUser(user.Name) // arg 0 is string of username
	if err != nil {
		return fmt.Errorf("error setting the user as '%v': %w", user.Name, err)
	}
	fmt.Printf("The user has been set to %v.\n", user.Name)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("arguments expected: 'username'")
	}
	username := cmd.Args[0]
	now := time.Now().UTC()
	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{CreatedAt: now, UpdatedAt: now, Name: username})
	if err != nil {
		return fmt.Errorf("error creating user=%v: %w", username, err)
	}
	s.S.SetUser(user.Name)
	fmt.Printf("User '%v' has been created.\n", user.Name)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("error, no arguments expected")
	}
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting all the registered users: %w", err)
	}
	fmt.Println("Deleted all registered users.")
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("error, no arguments expected")
	}
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}
	var out string
	for _, user := range users {
		name := user.Name + "\n"
		if user.Name == s.S.CurrentUsername {
			name = user.Name + " (current)\n"
		}
		out += name
	}
	fmt.Print(out)
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("error, expected args='time_betwee_reqs' in seconds")
	}
	seconds, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error converting '%v' to int: %w", cmd.Args[0], err)
	}
	ticker := time.NewTicker(time.Duration(seconds) * time.Second)
	fmt.Printf("Collecting feeds every %v seconds...\n\n", seconds)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

// called in HandlerAgg
func scrapeFeeds(s *State) error {
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting the next feed to fetch: %w", err)
	}
	now := time.Now().UTC()
	nowNullable := sql.NullTime{Time: now, Valid: true}

	err = s.Db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{ID: feed.ID, LastFetchedAt: nowNullable, UpdatedAt: now})
	if err != nil {
		return fmt.Errorf("error marking the fetched feed '%v': %w", feed.Url, err)
	}

	rss, err := gatorfeed.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed '%v': %w", feed.Url, err)
	}
	gatorfeed.CleanFeed(rss)
	fmt.Printf("Saving posts in the database for feed: %v; %v\n", feed.Name, feed.Url)
	for _, rssitem := range rss.Channel.Item {
		now := time.Now().UTC()
		_, err := s.Db.CreatePost(context.Background(), database.CreatePostParams{CreatedAt: now, UpdatedAt: now, Title: rssitem.Title,
			Url: rssitem.Link, Description: rssitem.Description, PublishedAt: rssitem.PubDate, FeedID: feed.ID})
		if err != nil {
			fmt.Println(fmt.Errorf("error creating post for '%v;%v': %w\n", rssitem.Title, rssitem.Link, err))
		}
	}
	return nil
}

func HandlerBrowse(s *State, c Command, user database.User) error {
	if len(c.Args) >= 2 {
		return fmt.Errorf("error, expected args='num_posts (optional)'")
	}
	var limit int
	var err error
	if len(c.Args) == 0 {
		limit = 2
	} else {
		limit, err = strconv.Atoi(c.Args[0])
		if err != nil {
			return fmt.Errorf("error, cannot convert '%v' to int: %w", c.Args[0], err)
		}
	}

	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{UserID: user.ID, Limit: int32(limit)})
	if err != nil {
		return fmt.Errorf("error getting posts for user '%v' with limit '%v': %w", user.Name, limit, err)
	}
	out := fmt.Sprintf("'%v' latest posts followed by '%v':\n", limit, user.Name)
	for _, post := range posts {
		out += "*************************************************************\n"
		out += fmt.Sprintf("Title: %v\nURL: %v\nFeed: %v\nPublished At: %v\nDescription: %v\n", post.Title, post.Url, post.FeedName, post.PublishedAt, post.Description)
	}
	fmt.Print(out)
	return nil
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("error, expecting args='name','url'")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	now := time.Now().UTC()
	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{CreatedAt: now, UpdatedAt: now, Name: name, Url: url, UserID: user.ID})
	if err != nil {
		return fmt.Errorf("error adding feed '%v' to the database: %w", name, err)
	}
	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{CreatedAt: now, UpdatedAt: now,
		UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("error '%v' tyring to follow '%v': %w", user.Name, feed.Url, err)
	}
	fmt.Printf("User '%v' has added feed: %v\n", user.Name, feed.Name)
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("error, no arguments expected")
	}
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting the feeds from the database: %w", err)
	}
	var out string
	for _, feed := range feeds {
		out += fmt.Sprintf("Name: %v\tURL: %v\tAdding User: %v\n", feed.Name, feed.Url, feed.AddingUser)
	}
	if out == "" {
		fmt.Println(out)
	} else {
		fmt.Print(out)
	}
	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("error, expected args='url'")
	}
	url := cmd.Args[0]
	feedToFollow, err := s.Db.GetFeedFromUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error getting feed to follow using url '%v': %w", url, err)
	}
	now := time.Now().UTC()
	followedData, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{CreatedAt: now, UpdatedAt: now,
		UserID: user.ID, FeedID: feedToFollow.ID})

	if err != nil {
		return fmt.Errorf("error '%v' could not follow '%v': %w", user.Name, url, err)
	}
	fmt.Printf("User '%v' has followed '%v'.\n", followedData.UserName, followedData.FeedName)
	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("error, expected args=none")
	}
	userFollowings, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting the followers for user '%v': %w", user.Name, err)
	}
	out := fmt.Sprintf("%v's followings:\n", user.Name)
	for _, userFollowing := range userFollowings {
		out += fmt.Sprintf("Feed: %v\tURL: %v\tAdding User: %v\n", userFollowing.FeedName, userFollowing.Url, userFollowing.AddingUser)
	}
	fmt.Println(out)
	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("error, expected args=url")
	}
	url := cmd.Args[0]
	feed, err := s.Db.GetFeedFromUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error getting feed from url '%v': %w", url, err)
	}
	err = s.Db.DeleteFeedFollowForUser(context.Background(), database.DeleteFeedFollowForUserParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("error deleting feedfollow=UserID: %v, FeedID: %v: %w", user.ID, feed.ID, err)
	}
	return nil
}

func HandlerHelp(s *State, cmd Command, c Commands) error {
	out := "Type 'bootdev-gator' before each command. All possible commands:\n"
	for name, description := range c.Descriptions {
		out += fmt.Sprintf("%v: %v\n", name, description)
	}
	fmt.Print(out)
	return nil
}
