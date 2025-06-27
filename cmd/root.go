package cmd

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	cobalias "github.com/spf13/cobra"
	"os"
	"strings"
)

var logLevel string

var rootCmd = &cobalias.Command{
	Use:   "kctl",
	Short: "kctl is a custom Kubernetes controller CLI",
	Long:  `"kctl" is a tool to test and run components of your custom Kubernetes controller`,
	PersistentPreRun: func(cmd *cobalias.Command, args []string) {
		level := parseLogLevel(logLevel)
		configureLogger(level)
	},
	Run: func(cmd *cobalias.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		log.Info().Msg("This is an info log")
		log.Debug().Msg("This is a debug log")
		log.Trace().Msg("This is a trace log")
		log.Warn().Msg("This is a warn log")
		log.Error().Msg("This is an error log")
		fmt.Println("")
	},
}

func configureLogger(level zerolog.Level) {
	// Set global time format and log level
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	zerolog.SetGlobalLevel(level)

	// Console writer setup
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: zerolog.TimeFieldFormat,
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.MessageFieldName,
		},
	}

	// Add caller for trace level
	if level == zerolog.TraceLevel {
		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			return fmt.Sprintf("%s:%d", file, line)
		}
		zerolog.CallerFieldName = "caller"

		consoleWriter.PartsOrder = []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.CallerFieldName,
			zerolog.MessageFieldName,
		}

		log.Logger = log.Output(consoleWriter).With().Caller().Logger()

	} else if level <= zerolog.DebugLevel {
		// For Debug or Info: human-readable console logs
		log.Logger = log.Output(consoleWriter)

	} else {
		// Raw JSON output (useful for log collection)
		log.Logger = log.Output(os.Stderr)
	}
}

func parseLogLevel(lvl string) zerolog.Level {
	switch strings.ToLower(lvl) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("Execution failed: ", err)
		log.Fatal().Msg(err.Error())
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Set log level: trace, debug, info, warn, error")
}
