// Copyright (c) 2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txauthor

import (
	"testing"

	"github.com/dcrlabs/ltcwallet/wallet/txrules"
	"github.com/dcrlabs/ltcwallet/wallet/txsizes"
	"github.com/ltcsuite/ltcd/ltcutil"
	"github.com/ltcsuite/ltcd/wire"
)

func p2pkhOutputs(amounts ...ltcutil.Amount) []*wire.TxOut {
	v := make([]*wire.TxOut, 0, len(amounts))
	for _, a := range amounts {
		outScript := make([]byte, txsizes.P2PKHOutputSize)
		v = append(v, wire.NewTxOut(int64(a), outScript))
	}
	return v
}

func makeInputSource(unspents []*wire.TxOut) InputSource {
	// Return outputs in order.
	currentTotal := ltcutil.Amount(0)
	currentInputs := make([]*wire.TxIn, 0, len(unspents))
	currentInputValues := make([]ltcutil.Amount, 0, len(unspents))
	f := func(target ltcutil.Amount) (ltcutil.Amount, []*wire.TxIn, []ltcutil.Amount, [][]byte, error) {
		for currentTotal < target && len(unspents) != 0 {
			u := unspents[0]
			unspents = unspents[1:]
			nextInput := wire.NewTxIn(&wire.OutPoint{}, nil, nil)
			currentTotal += ltcutil.Amount(u.Value)
			currentInputs = append(currentInputs, nextInput)
			currentInputValues = append(currentInputValues, ltcutil.Amount(u.Value))
		}
		return currentTotal, currentInputs, currentInputValues, make([][]byte, len(currentInputs)), nil
	}
	return InputSource(f)
}

func TestNewUnsignedTransaction(t *testing.T) {
	tests := []struct {
		UnspentOutputs   []*wire.TxOut
		Outputs          []*wire.TxOut
		RelayFee         ltcutil.Amount
		ChangeAmount     ltcutil.Amount
		InputSourceError bool
		InputCount       int
	}{
		0: {
			UnspentOutputs:   p2pkhOutputs(1e8),
			Outputs:          p2pkhOutputs(1e8),
			RelayFee:         1e3,
			InputSourceError: true,
		},
		1: {
			UnspentOutputs: p2pkhOutputs(1e8),
			Outputs:        p2pkhOutputs(1e6),
			RelayFee:       1e3,
			ChangeAmount: 1e8 - 1e6 - txrules.FeeForSerializeSize(1e3,
				txsizes.EstimateVirtualSize(1, 0, 0, 0, p2pkhOutputs(1e6), txsizes.P2WPKHPkScriptSize)),
			InputCount: 1,
		},
		2: {
			UnspentOutputs: p2pkhOutputs(1e8),
			Outputs:        p2pkhOutputs(1e6),
			RelayFee:       1e4,
			ChangeAmount: 1e8 - 1e6 - txrules.FeeForSerializeSize(1e4,
				txsizes.EstimateVirtualSize(1, 0, 0, 0, p2pkhOutputs(1e6), txsizes.P2WPKHPkScriptSize)),
			InputCount: 1,
		},
		3: {
			UnspentOutputs: p2pkhOutputs(1e8),
			Outputs:        p2pkhOutputs(1e6, 1e6, 1e6),
			RelayFee:       1e4,
			ChangeAmount: 1e8 - 3e6 - txrules.FeeForSerializeSize(1e4,
				txsizes.EstimateVirtualSize(1, 0, 0, 0, p2pkhOutputs(1e6, 1e6, 1e6), txsizes.P2WPKHPkScriptSize)),
			InputCount: 1,
		},
		4: {
			UnspentOutputs: p2pkhOutputs(1e8),
			Outputs:        p2pkhOutputs(1e6, 1e6, 1e6),
			RelayFee:       2.55e3,
			ChangeAmount: 1e8 - 3e6 - txrules.FeeForSerializeSize(2.55e3,
				txsizes.EstimateVirtualSize(1, 0, 0, 0, p2pkhOutputs(1e6, 1e6, 1e6), txsizes.P2WPKHPkScriptSize)),
			InputCount: 1,
		},

		// Test dust thresholds (546 for a 1e3 relay fee).
		5: {
			UnspentOutputs: p2pkhOutputs(1e8),
			Outputs: p2pkhOutputs(1e8 - 5459 - txrules.FeeForSerializeSize(1e3,
				txsizes.EstimateVirtualSize(1, 0, 0, 0, p2pkhOutputs(0), txsizes.P2WPKHPkScriptSize))),
			RelayFee:     1e3,
			ChangeAmount: 5459,
			InputCount:   1,
		},
		6: {
			UnspentOutputs: p2pkhOutputs(1e8),
			Outputs: p2pkhOutputs(1e8 - 5460 - txrules.FeeForSerializeSize(1e3,
				txsizes.EstimateVirtualSize(1, 0, 0, 0, p2pkhOutputs(0), txsizes.P2WPKHPkScriptSize))),
			RelayFee:     1e3,
			ChangeAmount: 5460,
			InputCount:   1,
		},

		// Test dust thresholds (1392.3 for a 2.55e3 relay fee).
		7: {
			UnspentOutputs: p2pkhOutputs(1e8),
			Outputs: p2pkhOutputs(1e8 - 13920 - txrules.FeeForSerializeSize(2.55e3,
				txsizes.EstimateVirtualSize(1, 0, 0, 0, p2pkhOutputs(0), txsizes.P2WPKHPkScriptSize))),
			RelayFee:     2.55e3,
			ChangeAmount: 13920,
			InputCount:   1,
		},
		8: {
			UnspentOutputs: p2pkhOutputs(1e8),
			Outputs: p2pkhOutputs(1e8 - 13930 - txrules.FeeForSerializeSize(2.55e3,
				txsizes.EstimateVirtualSize(1, 0, 0, 0, p2pkhOutputs(0), txsizes.P2WPKHPkScriptSize))),
			RelayFee:     2.55e3,
			ChangeAmount: 13930,
			InputCount:   1,
		},

		// Test two unspent outputs available but only one needed
		// (tested fee only includes one input rather than using a
		// serialize size for each).
		9: {
			UnspentOutputs: p2pkhOutputs(1e8, 1e8),
			Outputs: p2pkhOutputs(1e8 - 5460 - txrules.FeeForSerializeSize(1e3,
				txsizes.EstimateVirtualSize(1, 0, 0, 0, p2pkhOutputs(0), txsizes.P2WPKHPkScriptSize))),
			RelayFee:     1e3,
			ChangeAmount: 5460,
			InputCount:   1,
		},

		// Test that second output is not included to make the change
		// output not dust and be included in the transaction.
		//
		// It's debatable whether or not this is a good idea, but it's
		// how the function was written, so test it anyways.
		10: {
			UnspentOutputs: p2pkhOutputs(1e8, 1e8),
			Outputs: p2pkhOutputs(1e8 - 5450 - txrules.FeeForSerializeSize(1e3,
				txsizes.EstimateVirtualSize(1, 0, 0, 0, p2pkhOutputs(0), txsizes.P2WPKHPkScriptSize))),
			RelayFee:     1e3,
			ChangeAmount: 5450,
			InputCount:   1,
		},

		// Test two unspent outputs available where both are needed.
		11: {
			UnspentOutputs: p2pkhOutputs(1e8, 1e8),
			Outputs:        p2pkhOutputs(1e8),
			RelayFee:       1e3,
			ChangeAmount: 1e8 - txrules.FeeForSerializeSize(1e3,
				txsizes.EstimateVirtualSize(2, 0, 0, 0, p2pkhOutputs(1e8), txsizes.P2WPKHPkScriptSize)),
			InputCount: 2,
		},

		// Test that zero change outputs are not included
		// (ChangeAmount=0 means don't include any change output).
		12: {
			UnspentOutputs: p2pkhOutputs(1e8),
			Outputs:        p2pkhOutputs(1e8),
			RelayFee:       0,
			ChangeAmount:   0,
			InputCount:     1,
		},
	}

	changeSource := &ChangeSource{
		NewScript: func() ([]byte, error) {
			// Only length matters for these tests.
			return make([]byte, txsizes.P2WPKHPkScriptSize), nil
		},
		ScriptSize: txsizes.P2WPKHPkScriptSize,
	}

	for i, test := range tests {
		inputSource := makeInputSource(test.UnspentOutputs)
		tx, err := NewUnsignedTransaction(test.Outputs, test.RelayFee, inputSource, changeSource)
		switch e := err.(type) {
		case nil:
		case InputSourceError:
			if !test.InputSourceError {
				t.Errorf("Test %d: Returned InputSourceError but expected "+
					"change output with amount %v", i, test.ChangeAmount)
			}
			continue
		default:
			t.Errorf("Test %d: Unexpected error: %v", i, e)
			continue
		}
		if tx.ChangeIndex < 0 {
			if test.ChangeAmount != 0 {
				t.Errorf("Test %d: No change output added but expected output with amount %v",
					i, test.ChangeAmount)
				continue
			}
		} else {
			changeAmount := ltcutil.Amount(tx.Tx.TxOut[tx.ChangeIndex].Value)
			if test.ChangeAmount == 0 {
				t.Errorf("Test %d: Included change output with value %v but expected no change",
					i, changeAmount)
				continue
			}
			if changeAmount != test.ChangeAmount {
				t.Errorf("Test %d: Got change amount %v, Expected %v",
					i, changeAmount, test.ChangeAmount)
				continue
			}
		}
		if len(tx.Tx.TxIn) != test.InputCount {
			t.Errorf("Test %d: Used %d outputs from input source, Expected %d",
				i, len(tx.Tx.TxIn), test.InputCount)
		}
	}
}
