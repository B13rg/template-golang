package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ConfigureLogger(debug bool) {
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		switch RootConfig.LogLevel {
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "info":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "warn":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case "error":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		case "fatal":
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
		case "panic":
			zerolog.SetGlobalLevel(zerolog.PanicLevel)
		default:
			log.Fatal().Msg("invalid log level: " + RootConfig.LogLevel)
		}
	}

	log.Logger = SetupLogger(RootConfig.Color)
}

// Configure zerolog with some defaults and cleanup error formatting.
func SetupLogger(enableColor bool) zerolog.Logger {
	//nolint:exhaustruct
	consoleWriter := zerolog.ConsoleWriter{
		Out:     os.Stderr,
		NoColor: !enableColor,
		FormatErrFieldValue: func(err interface{}) string {
			// https://github.com/rs/zerolog/blob/a21d6107dcda23e36bc5cfd00ce8fdbe8f3ddc23/console.go#L21
			colorRed := 31
			colorBold := 1
			s := strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(err.(string), "\\t", " "), "\\n", " | ",
				), "|  |", "|")

			return Colorize(Colorize(s, colorBold, !enableColor), colorRed, !enableColor)
		},
		// Other fields:
		// TimeFormat, TimeLocation, PartsOrder, PartsExclude,
		// FieldsOrder, FieldsExclude, FormatTimestamp, FormatLevel,
		// FormatCaller, FormatMessage, FormatFieldName, FormatFieldValue,
		// FormatErrFieldName, FormatExtra, FormatPrepare
	}

	return log.Output(consoleWriter)
}

// Colorize function from zerolog console.go file to replicate their coloring functionality.
// Source: https://github.com/rs/zerolog/blob/a21d6107dcda23e36bc5cfd00ce8fdbe8f3ddc23/console.go#L389
// Replicated here because it's a private function.
func Colorize(input interface{}, colorNum int, disabled bool) string {
	e := os.Getenv("NO_COLOR")

	if disabled || (e != "" || colorNum == 0) {
		// escape hatch for disabled coloring or empty input
		return fmt.Sprintf("%s", input)
	}

	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", colorNum, input)
}
