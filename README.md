# Gator RSS Aggregator

Gator is a command-line RSS aggregator written in Go. It lets you subscribe to RSS feeds, fetches new posts on a schedule, and lets you browse your personalized feed collection—all from your terminal.

---

## Features

- Register and manage users
- Add and follow RSS feeds
- Periodically fetch and store new posts from feeds
- Browse posts from feeds you follow
- List all feeds and your subscriptions
- Backed by PostgreSQL for reliable storage

---

## Requirements

- **Go** (version 1.21 or newer recommended)
- **PostgreSQL** (version 12+ recommended)

---

## Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/mcoluomo/RSS-Aggregator.git
   cd RSS-Aggregator
2. **Install the CLI:**

   Make sure `$GOPATH/bin` is in your `PATH`, then run:

   ```sh
   go install ./...
   ```

   This will build and install the `gator` CLI binary.

   ```sh
   chmod +x gator
   sudo mv gator /usr/local/bin/
   ```

   This will put the statically linked binary on your path

   You can check that `gator` is in your `$PATH`:

   ```sh
   which gator
   ```

   You should see `/usr/local/bin/gator` or similar.

---

## Database Setup

1. **Start PostgreSQL** and create a database (e.g., `gator`):

   ```sh
   createdb gator
   ```

2. **Run the migrations** (using [goose](https://github.com/pressly/goose) or your preferred tool):

   ```sh
   goose -dir internal/sql/schema postgres "postgres://postgres:<yourpassword>@localhost:5432/gator?sslmode=disable" up
   ```

---

## Configuration

Create a config file at `~/.gatorconfig.json` with the following content:

```json
{
  "db_url": "postgres://postgres:<yourpassword>@localhost:5432/gator?sslmode=disable",
  "current_user_name": "[None]"
}
```

- Replace `<yourpassword>` with your actual Postgres password.

---

## Usage

Run the CLI with:

```sh
gator <command> [arguments]
```

### Common Commands

- **Register a new user:**
  ```sh
  gator register <username>
  ```

- **Login as a user:**
  ```sh
  gator login <username>
  ```

- **Add a new feed and follow it:**
  ```sh
  gator addfeed "<feed name>" <feed url>
  ```

- **Follow an existing feed:**
  ```sh
  gator follow <feed url>
  ```

- **List all feeds:**
  ```sh
  gator feeds
  ```

- **List feeds you are following:**
  ```sh
  gator following
  ```

- **Fetch new posts from feeds (every 10 minutes):**
  ```sh
  gator agg 10m
  ```

- **Browse your latest posts:**
  ```sh
  gator browse 5
  ```

- **Reset all users (dangerous!):**
  ```sh
  gator reset
  ```

---

## How It Works

- **Users** register and log in.
- **Feeds** are added to the system and can be followed by any user.
- The **aggregator** fetches new posts from feeds on a schedule and stores them in the database.
- **Users** can browse posts from the feeds they follow.

---

## Development

- Code is organized in `internal/cli`, `internal/database`, and `internal/config`.
- SQL queries and migrations are in `internal/sql/`.
- Uses [sqlc](https://sqlc.dev/) for type-safe database access.

---

## License

MIT License – see [LICENSE](./LICENSE) file for details.
