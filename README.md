jWalterWeatherman
=================

A simple printing and logging library for go

![Always Leave A Note](http://spf13.github.com/jwalterweatherman/and_that__s_why_you_always_leave_a_note_by_jonnyetc-d57q7um.jpg)
Graphic by [JonnyEtc](http://jonnyetc.deviantart.com/art/And-That-s-Why-You-Always-Leave-a-Note-315311422)

JWW is primarily a convenience wrapper around the
excellent standard log library.

I really wanted a very straightforward library that could seamlessly do
the following things.

1. Replace all the println, printf, etc statements thought my code with
   something more useful
2. Allow the user to easily control what levels are printed
3. Allow the user to easily control what levels are logged
4. Provide good feedback (which can easily be logged) to the user
5. Due to 2 & 3 provide easy verbose mode for output and logs
6. Not have any unnecessary initialization cruft. Just use it.

# Usage

## Step 1. Use it
Put calls throughout your source based on type of feedback.

Available Loggers are:

 * TRACE
 * DEBUG
 * INFO
 * WARN
 * ERROR
 * CRITICAL
 * FATAL

These each are loggers based on the log standard library and follow the
standard usage. Eg..

```go
    import (
        jww "github.com/spf13/jwalterweatherman"
    )

    ...

    if err != nil {
        jww.ERROR.Println(err)
    }

    // this isn’t that important, but they may want to know
    jww.INFO.Printf("information %q", response)

```


## Step 2. Optionally configure it

By default:
 * Debug, Trace & Info goto /dev/null
 * Warn and above is logged (when a log file/io.Writer is provided)
 * Error and above is printed to the terminal (stdout)


### Changing the thresholds

The threshold can be changed at any time, but will only affect calls that
execute after the change was made.

This is very useful if your application has a verbose mode. Of course you
can decide what verbose means to you or even have multiple levels of
verbosity.


```go
    import (
        jww "github.com/spf13/jwalterweatherman"
    )

    if Verbose {
        jww.SetLogThreshold(jww.LevelTrace)
        jww.SetOutputThreshold(jww.LevelInfo)
    }
```

### Using a temp log file

JWW conveniently creates a temporary file and sets the log Handle to
a io.Writer created for it. You should call this early in your application
initialization routine as it will only log calls made after it is executed. 

```go
    import (
        jww "github.com/spf13/jwalterweatherman"
    )

    jww.UseTempLogFile("YourAppName") 

```

### Setting a log file

JWW can log to any file you provide a path to (provided it’s writable).
Will only append to this file.


```go
    import (
        jww "github.com/spf13/jwalterweatherman"
    )

    jww.SetLogFile("/path/to/logfile") 

```


# More information

This is an early release. I’ve been using it for a while and this is the
third interface I’ve tried. I like this one pretty well, but no guarantees
that it won’t change a bit.

I wrote this for use in [hugo](http://hugo.spf13.com). If you are looking
for a static website engine that’s super fast please checkout Hugo.
