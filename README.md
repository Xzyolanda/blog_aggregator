# Gator RSS Aggregator

Gator is a command-line RSS aggregator written in Go. It stores users, feeds, follows, and posts in PostgreSQL. The `agg` command continuously fetches feeds; `browse` shows recent posts from feeds followed by the current user.

## Requirements

- Go 1.27.2 or newer
- PostgreSQL, running and reachable from this machine
- The Goose migration CLI, to create or update the database schema

## Install the CLI

From the repository root, install the current checkout with:

```sh
go install .
```

To install the published module instead, run:

```sh
go install github.com/Xzyolanda/blog_aggregator@latest
```

Go installs executables to `GOBIN`, or to `$(go env GOPATH)/bin` when `GOBIN` is unset. Make sure that directory is on your `PATH`. The module currently installs as `blog_aggregator`, so the examples below use that command name.

## Configure PostgreSQL

Create a database for the aggregator. For a local PostgreSQL installation, for example:

```sh
createdb gator
```

Use a PostgreSQL connection URL like this, replacing the username and password with your local credentials:

```text
postgres://postgres:your_password@localhost:5432/gator?sslmode=disable
```

The program reads its configuration from `~/.gatorconfig.json`. Create that file with the same connection URL:

```json
{
  "db_url": "postgres://postgres:your_password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Registering or logging in sets `current_user_name` automatically.

## Create or update the database schema

Install Goose if it is not already available:

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
```

From the repository root, apply all migrations using the same connection URL configured above:

```sh
goose -dir ./sql/schema postgres "postgres://postgres:your_password@localhost:5432/gator?sslmode=disable" up
```

The migrations create the users, feeds, follows, fetch timestamps, and posts tables. Run this command again after pulling changes that add migrations.

## Run the program

Register a user; this also sets that user as the current user:

```sh
blog_aggregator register alice
```

Add a feed. Adding a feed also follows it for the current user:

```sh
blog_aggregator addfeed "Hacker News" "https://hnrss.org/newest"
blog_aggregator addfeed "Lane's Blog" "https://www.wagslane.dev/index.xml"
```

Run the aggregator in a separate terminal. It fetches feeds continuously, with the interval supplied as a Go duration. Choose a reasonable interval to avoid making requests too frequently:

```sh
blog_aggregator agg 1m
```

Browse the two most recent posts by default, or provide a positive limit:

```sh
blog_aggregator browse
blog_aggregator browse 10
```

## Other commands

- `blog_aggregator feeds` — list known feeds.
- `blog_aggregator follow <url>` — follow an existing feed.
- `blog_aggregator following` — list feeds followed by the current user.
- `blog_aggregator unfollow <url>` — stop following a feed.
- `blog_aggregator login <name>` — switch the current user.
- `blog_aggregator users` — list users and identify the current user.
- `blog_aggregator reset` — delete all users. This is destructive; related feeds and follows are removed as well.

You can run commands from the repository without installing the binary by replacing `blog_aggregator` with `go run .`, for example `go run . browse`.
