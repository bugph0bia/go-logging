# go-logging

Simple logging & log rotation using slog and lumberjack

## Usage

### Simple to use

```go
package main

import "github.com/bugph0bia/go-logging"

func main() {
    // Create logger with default values
    logger := logging.NewLogger("log.txt")

    // output logs
    logger.Debug("message")
    logger.Info("message", "attr1", 10)
    logger.Warn("message", "attr1", 10, "attr2", 20)
    logger.Error("message")
}

// The log output is follow:
//
// 2024/10/24 11:22:33 DEBUG: message
// 2024/10/24 11:22:33 INFO [attr1=10]: message
// 2024/10/24 11:22:33 WARN [attr1=10, attr2=20]: message
// 2024/10/24 11:22:33 ERROR: message
```

### If you don't want to pass loggers between functions

```go
package main

import (
    "log/slog"

    "github.com/bugph0bia/go-logging"
)

func main() {
    // Create logger and give to slog.
    slog.SetDefault(logging.NewLogger("log.txt"))

    // output logs
    slog.Debug("message")
}
```

### Options

To change an option, first get the handler with `Logging.NewHandler`, change options, and then call `logging.NewLoggerFromHandler` to create the logger.  
The sample code below is an example of the options that can be changed and their default values.  

```go
func main() {
    // Create handler
    handler := logging.NewHandler("log.txt")


    //////// Change log output behavior ////////

    // Output log level
    handler.Option.Level = slog.LevelInfo

    // Flag to output log to stdout
    handler.WithStdout = true

    //////// Change log format ////////

    // Format of log line (Use tag string)
    handler.Format.Line = fmt.Sprintf("${Datetime} ${Level} ${Attrs}: ${Message}")

    // Format of log line (Use const values)
    handler.Format.Line = fmt.Sprintf("%s %s %s: %s", logging.FDatetime, logging.FLevel, logging.FAttrs, logging.FMessage)

    // Format of datetime
    handler.Format.Datetime = "2006/01/02 15:04:05"

    // Character between the key and value of the attribute
    handler.Format.AttrBetween = "="

    // Separator between attributes
    handler.Format.AttrDelimiter = ", "

    // Attribute list prefix
    handler.Format.AttrPrefix = "["

    // Attribute list suffix
    handler.Format.AttrSuffix = "]"

    //////// Change log rotation ////////

    // Max size of a log file (MB)
    // Rotate beyond this size
    handler.RotateLogger.MaxSizeMB = 1

    // Maximum number of backup files
    // If zero, no limit
    handler.RotateLogger.MaxBackups = 10

    // Maximum number of days to keep backup files
    // If zero, no limit
    handler.RotateLogger.MaxAge = 0

    // Flag to set backup file time to local time
    handler.RotateLogger.LocalTime = true

    // Flag to gzip compress backup files
    handler.RotateLogger.Compress = false


    // Create Logger
    logger := logging.NewLoggerFromHandler(handler)
}
```

## Thanks

- [lumberjack](https://github.com/natefinch/lumberjack)
