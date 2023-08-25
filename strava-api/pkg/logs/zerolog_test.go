package logs_test

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/rzajac/zltest"

	"yellow-jersey/pkg/logs"
)

func Test_Logger_New(t *testing.T) {
	// Crate zerolog test helper.
	tst := zltest.New(t)

	logs.New(logs.WithWriter(tst))
	logs.Logger.Info().Msg("test global logger")

	// Test if log messages were generated properly.
	ent := tst.LastEntry()
	ent.ExpMsg("test global logger")
	ent.ExpLevel(zerolog.InfoLevel)
}

func Test_Logger_Error(t *testing.T) {
	// Crate zerolog test helper.
	tst := zltest.New(t)

	logs.New(logs.WithWriter(tst))
	logs.Logger.Error().Msg("huge error")

	// Test if log messages were generated properly.
	ent := tst.LastEntry()
	ent.ExpMsg("huge error")
	ent.ExpLevel(zerolog.ErrorLevel)
}

func Test_Logger_MetaData(t *testing.T) {
	// Crate zerolog test helper.
	tst := zltest.New(t)

	logs.New(logs.WithWriter(tst), logs.WithVersion("0.1.0"), logs.WithService("test-service"))
	logs.Logger.Info().Msg("test global logger")

	// Test if log messages were generated properly.
	ent := tst.LastEntry()
	ent.ExpStr("version", "0.1.0")
	ent.ExpStr("service", "test-service")
	ent.ExpMsg("test global logger")
	ent.ExpLevel(zerolog.InfoLevel)
}

func Test_Logger_Debug(t *testing.T) {
	// Crate zerolog test helper.
	tst := zltest.New(t)

	logs.New(logs.WithWriter(tst), logs.WithVersion("0.1.0"), logs.WithService("test-service"), logs.WithDebug())
	logs.Logger.Debug().Msg("test global logger")

	// Test if log messages were generated properly.
	ent := tst.LastEntry()
	ent.ExpStr("version", "0.1.0")
	ent.ExpStr("service", "test-service")
	ent.ExpMsg("test global logger")
	ent.ExpLevel(zerolog.DebugLevel)
}
