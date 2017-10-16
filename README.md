=== Usage

Use the parse command to parse a crontab entry:

`go run main.go parse "*/15 0 1,15 * 1-5 /usr/bin/find"`

Any invalid entry will panic.

TODO:

- [] Support Sunday as 0 and 7 in 'day of week'
- [] Support 'extended' values like @daily, @weekly.
