package gg

import (
	"fmt"
	"os"
	"time"
)

/*
Shortcut for creating `LogTime` with the current time and with a message
generated from the inputs via `Str`.
*/
func LogTimeNow(msg ...any) LogTime {
	return LogTime{Start: time.Now(), Msg: Str(msg...)}
}

/*
Shortcut for logging execution timing to stderr. Usage examples:

	defer gg.LogTimeNow(`some_activity`).LogStart().LogEnd()
	// perform some activity

	defer gg.LogTimeNow(`some_activity`).LogEnd()
	// perform some activity

	timer := gg.LogTimeNow(`some_activity`).LogStart()
	// perform some activity
	timer.LogEnd()

	timer := gg.LogTimeNow(`some_activity`)
	// perform some activity
	timer.LogEnd()
*/
type LogTime struct {
	Start time.Time
	Msg   string
}

/*
Logs the beginning of the activity denoted by `.Msg`:

	[some_activity] starting

Note that when logging start AND end, the time spent in `.LogStart` is
unavoidably included into the difference. If you want precise timing,
avoid logging the start.
*/
func (self LogTime) LogStart() LogTime {
	fmt.Fprintf(os.Stderr, "[%v] starting\n", self.Msg)
	return self
}

/*
Prints the end of the activity denoted by `.Msg`, with time elapsed since the
beginning:

	[some_activity] done in <duration>

If deferred, this will detect the current panic, if any, and print the following
instead:

	[some_activity] failed in <duration>
*/
func (self LogTime) LogEnd() LogTime {
	since := time.Since(self.Start)
	err := AnyErrTracedAt(recover(), 1)

	if err != nil {
		fmt.Fprintf(os.Stderr, "[%v] failed in %v\n", self.Msg, since)
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "[%v] done in %v\n", self.Msg, since)
	return self
}
