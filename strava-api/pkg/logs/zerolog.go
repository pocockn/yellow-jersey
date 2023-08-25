package logs

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

// Logger is the global zerolog logger which we can customise in our constructors.
// Usage - logs.Logger.Info().Msgf("%d attempt at connecting to the DB", i)
var Logger zerolog.Logger

// New sets the global log level and assigns the global logger to the package variable.
// It expects the version number and service name to attach additional config to the logs.
// We set the default writer as os.Stdout. This can be customized with the WithWriter function option.
// logs.New(logs.WithWriter(os.Stderr))
func New(opts ...NewFuncOption) {
	logLevel := zerolog.InfoLevel
	zerolog.SetGlobalLevel(logLevel)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	for _, o := range opts {
		logger = o(&logger)
	}

	Logger = logger
}

// NewFuncOption is a functional option for the New function.
type NewFuncOption func(logger *zerolog.Logger) zerolog.Logger

// WithDebug sets the log level to debug. Used in all environments except production.
func WithDebug() NewFuncOption {
	return func(zl *zerolog.Logger) zerolog.Logger {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		return *zl
	}
}

// WithVersion sets the version of the service in the logs.
func WithVersion(version string) NewFuncOption {
	return func(zl *zerolog.Logger) zerolog.Logger {
		return zl.With().Str("version", version).Logger()
	}
}

// WithService adds the service name to the logs.
func WithService(service string) NewFuncOption {
	return func(zl *zerolog.Logger) zerolog.Logger {
		return zl.With().Str("service", service).Logger()
	}
}

// WithWriter specifies where to write the logs.
func WithWriter(writer io.Writer) NewFuncOption {
	return func(zl *zerolog.Logger) zerolog.Logger {
		return zerolog.New(writer).With().Timestamp().Logger()
	}
}
