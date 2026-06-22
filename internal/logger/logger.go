package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var global zerolog.Logger

func init() {
	Setup("info", false)
}

func Setup(level string, jsonFormat bool) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)

	if jsonFormat {
		global = zerolog.New(os.Stdout).
			With().
			Timestamp().
			Caller().
			Logger()
	} else {
		global = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}).
			With().
			Timestamp().
			Caller().
			Logger()
	}
}

func Debug() *zerolog.Event { return global.Debug() }
func Info() *zerolog.Event  { return global.Info() }
func Warn() *zerolog.Event  { return global.Warn() }
func Error() *zerolog.Event { return global.Error() }
func Fatal() *zerolog.Event { return global.Fatal() }
func Get() zerolog.Logger   { return global }
