package spv

import (
	"github.com/dcrlabs/ltcwallet/spv/blockntfns"
	"github.com/dcrlabs/ltcwallet/spv/chanutils"
	"github.com/dcrlabs/ltcwallet/spv/filterdb"
	"github.com/dcrlabs/ltcwallet/spv/pushtx"
	"github.com/dcrlabs/ltcwallet/spv/query"

	"github.com/btcsuite/btclog"
	"github.com/ltcsuite/ltcd/addrmgr"
	"github.com/ltcsuite/ltcd/blockchain"
	"github.com/ltcsuite/ltcd/connmgr"
	"github.com/ltcsuite/ltcd/peer"
	"github.com/ltcsuite/ltcd/txscript"
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
	log = btclog.Disabled
}

// UseLogger uses a specified Logger to output package logging info.
// This should be used in preference to SetLogWriter if the caller is also
// using btclog.
func UseLogger(logger btclog.Logger) {
	log = logger
	blockchain.UseLogger(logger)
	txscript.UseLogger(logger)
	peer.UseLogger(logger)
	addrmgr.UseLogger(logger)
	blockntfns.UseLogger(logger)
	pushtx.UseLogger(logger)
	connmgr.UseLogger(logger)
	query.UseLogger(logger)
	filterdb.UseLogger(logger)
	chanutils.UseLogger(logger)
}

// LogGenerate is an interface capable of generating loggers.
type LogGenerator interface {
	NewLogger(name string) btclog.Logger
}

// UseLogGenerator will set the logger for all subpackages with unique names.
func UseLogGenerator(g LogGenerator) {
	log = g.NewLogger("SPV")
	blockchain.UseLogger(g.NewLogger("BLKCHN"))
	txscript.UseLogger(g.NewLogger("TXSCRIPT"))
	peer.UseLogger(g.NewLogger("PEER"))
	addrmgr.UseLogger(g.NewLogger("ADRMGR"))
	blockntfns.UseLogger(g.NewLogger("BLKNOTE"))
	pushtx.UseLogger(g.NewLogger("PUSHTX"))
	connmgr.UseLogger(g.NewLogger("CONNMGR"))
	query.UseLogger(g.NewLogger("QUERY"))
}
