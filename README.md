# gator

This project was created for the boot.dev course. Gator is a basic RSS feed aggregator written in Go using a Postgres database.

# Install

You'll need to have Postgres and Go installed on your system.

Install gator using `go install github.com/Rodabaugh/gator@latest`

In your home directory, create a config file called ".gatorconfig.json"

It should contain something like this:

`{"db_url":"postgres://(YOUR POSTGRES DATABASE)"}`

Typically:

`{"db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"}`

# Usage

Gator is used by with `gator command (arguments)`

When you first install gator, create a user for yourself with `gator register <username>`

Other users can register the same way, and you can change who is logged in with the login command. `gator login <username>`

Feeds can be added with `gator addfeed <name> <url>`

Feeds can be listed with `gator feeds`

You can follow feeds added by other users with `gator follow <url>`

Unfollowing works as expected with `gator unfollow <url>`

To aggregate posts from the feeds, run `gator agg <timeBetweenRequests> (e.g. 30s, 1m, 1h)` This can be done in a second terminal as a kind of service while you use the program in another terminal.

To browse aggregated posts, use `gator browse` This defaults to 2 posts, but you can add an additional argument to specify how many posts you would like. e.g. `gator browse 24`

Additional command can be found with `gator help`

The database can be fully reset (this is a **irreversible** and **destructive** process). This can be done with `gator reset`