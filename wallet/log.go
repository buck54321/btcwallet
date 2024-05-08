// Copyright (c) 2015 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wallet

import (
	"github.com/btcsuite/btclog"
	"github.com/dcrlabs/ltcwallet/chain"
	"github.com/dcrlabs/ltcwallet/spv"
	"github.com/dcrlabs/ltcwallet/waddrmgr"
	"github.com/dcrlabs/ltcwallet/walletdb/migration"
	"github.com/dcrlabs/ltcwallet/wtxmgr"
)

// log is a logger that is initialized with no output filters.  This
// means the package will not perform any logging by default until the caller
// requests it.
var log btclog.Logger

// The default amount of logging is none.
func init() {
	DisableLog()
}

// DisableLog disables all library log output.  Logging output is disabled
// by default until either UseLogger or SetLogWriter are called.
func DisableLog() {
	UseLogger(btclog.Disabled)
}

// UseLogger uses a specified Logger to output package logging info.
// This should be used in preference to SetLogWriter if the caller is also
// using btclog.
func UseLogger(logger btclog.Logger) {
	log = logger

	migration.UseLogger(logger)
	waddrmgr.UseLogger(logger)
	wtxmgr.UseLogger(logger)
	chain.UseLogger(logger)
}

// LogGenerate is an interface capable of generating loggers.
type LogGenerator interface {
	NewLogger(name string) btclog.Logger
}

// UseLogGenerator will set the logger for all subpackages with unique names.
func UseLogGenerator(g LogGenerator) {
	log = g.NewLogger("WLLT")
	migration.UseLogger(g.NewLogger("MGRTN"))
	waddrmgr.UseLogger(g.NewLogger("ADRMGR"))
	wtxmgr.UseLogger(g.NewLogger("TXMGR"))
	chain.UseLogger(g.NewLogger("CHAIN"))
	spv.UseLogGenerator(g)
}

// pickNoun returns the singular or plural form of a noun depending
// on the count n.
func pickNoun(n int, singular, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}

// LogClosure is a closure that can be printed with %v to be used to
// generate expensive-to-create data for a detailed log level and avoid doing
// the work if the data isn't printed.
type logClosure func() string

// String invokes the log closure and returns the results string.
func (c logClosure) String() string {
	return c()
}

// newLogClosure returns a new closure over the passed function which allows
// it to be used as a parameter in a logging function that is only invoked when
// the logging level is such that the message will actually be logged.
func newLogClosure(c func() string) logClosure {
	return logClosure(c)
}
