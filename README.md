[![GoDoc](http://godoc.org/github.com/Willyham/tokei?status.png)](http://godoc.org/github.com/Willyham/tokei) 

=== Usage

Tokei is a simple and fast library for parsing and scheduling cron tasks. It can tell you when a cron will fire
at any point in the future, or it can give you a timer which fires every time the cron does.

```golang
expression, err := tokei.Parse("*/10 * * * *")
if err != nil {
  // handle err
}
schedule := tokei.NewScheduleUTC(expression)

// Get the next time that matches the cron
schedule.Next()

// Get the next 5 times which match the cron
schedule.Project(5)

// Get a timer which fires when the cron matches:
timer := schedule.Timer()

go func() {
  for {
    fired := <-timer.Next()
    fmt.Println("Fired at", fired)
  }
}()

go timer.Start()
```

==== Improvements

Tokei only currently supports standard cron entries rather than extended ones.

TODO:

- [] Support Sunday as 0 and 7 in 'day of week'
- [] Support 'extended' values like @daily, @weekly.
