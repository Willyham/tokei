[![GoDoc](http://godoc.org/github.com/Willyham/tokei?status.png)](http://godoc.org/github.com/Willyham/tokei) [![Go Report Card](https://goreportcard.com/badge/github.com/Willyham/tokei)](https://goreportcard.com/report/github.com/Willyham/tokei)

## Usage

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

### Benchmarks

Tokei is pretty quick, not that speed should be an issue for the kinds of things it can be used for. Nevertheless,
no golang library will be considered usable by the community without some arbitrary benchmarks showing how quick it is.

These benchmarks are generating the next firing time for some common cases; "years in future" is 'worst' case scenario where the first matching
entry doesn't occur until about 4 years from the start time.

```
BenchmarkNext/all-4                      	 5000000	       302 ns/op	      80 B/op	       2 allocs/op
BenchmarkNext/every_10_minutes-4         	 5000000	       311 ns/op	      80 B/op	       2 allocs/op
BenchmarkNext/waking_hours-4             	 3000000	       562 ns/op	      80 B/op	       2 allocs/op
BenchmarkNext/years_in_future-4          	   50000	     23230 ns/op	     272 B/op	       6 allocs/op
```

And calculating the next 5 firing times for the same entries:

```
BenchmarkProject/all-4                   	 1000000	      1495 ns/op	     368 B/op	       6 allocs/op
BenchmarkProject/every_10_minutes-4      	  500000	      2508 ns/op	     368 B/op	       6 allocs/op
BenchmarkProject/waking_hours-4          	  200000	      6516 ns/op	     368 B/op	       6 allocs/op
BenchmarkProject/years_in_future-4       	   50000	     25097 ns/op	     560 B/op	      10 allocs/op
```

### Improvements

Tokei only currently supports standard cron entries rather than extended ones.

TODO:

- [] Support Sunday as 0 and 7 in 'day of week'
- [] Support 'extended' values like @daily, @weekly.
