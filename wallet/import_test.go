package wallet

import (
	"encoding/binary"
	"strings"
	"testing"

	"github.com/dcrlabs/ltcwallet/waddrmgr"
	"github.com/ltcsuite/ltcd/chaincfg"
	"github.com/ltcsuite/ltcd/ltcutil"
	"github.com/ltcsuite/ltcd/ltcutil/hdkeychain"
	"github.com/ltcsuite/ltcd/txscript"
	"github.com/stretchr/testify/require"
)

func hardenedKey(key uint32) uint32 {
	return key + hdkeychain.HardenedKeyStart
}

func deriveAcctPubKey(t *testing.T, root *hdkeychain.ExtendedKey,
	scope waddrmgr.KeyScope, paths ...uint32) *hdkeychain.ExtendedKey {

	path := []uint32{hardenedKey(scope.Purpose), hardenedKey(scope.Coin)}
	path = append(path, paths...)

	var (
		currentKey = root
		err        error
	)
	for _, pathPart := range path {
		currentKey, err = currentKey.Derive(pathPart)
		require.NoError(t, err)
	}

	// The Neuter() method checks the version and doesn't know any
	// non-standard methods. We need to convert them to standard, neuter,
	// then convert them back with the target extended public key version.
	pubVersionBytes := make([]byte, 4)
	copy(pubVersionBytes, chaincfg.TestNet4Params.HDPublicKeyID[:])
	switch {
	case strings.HasPrefix(root.String(), "uprv"):
		binary.BigEndian.PutUint32(pubVersionBytes, uint32(
			waddrmgr.HDVersionTestNetBIP0049,
		))

	case strings.HasPrefix(root.String(), "vprv"):
		binary.BigEndian.PutUint32(pubVersionBytes, uint32(
			waddrmgr.HDVersionTestNetBIP0084,
		))
	}

	currentKey, err = currentKey.CloneWithVersion(
		chaincfg.TestNet4Params.HDPrivateKeyID[:],
	)
	require.NoError(t, err)
	currentKey, err = currentKey.Neuter()
	require.NoError(t, err)
	currentKey, err = currentKey.CloneWithVersion(pubVersionBytes)
	require.NoError(t, err)

	return currentKey
}

type testCase struct {
	name               string
	masterPriv         string
	accountIndex       uint32
	addrType           waddrmgr.AddressType
	expectedScope      waddrmgr.KeyScope
	expectedAddr       string
	expectedChangeAddr string
}

// TODO: The test cases below use the correct key scopes with ltc coin id = 2.
// When ltcwallet first forked, they neglected to change the coin id from the
// bitcoin value = 0. Check to make sure that imported wallets that may have
// created addresses with the incorrect coin ids are still imported correctly.

var (
	testCases = []*testCase{{
		name: "bip44 with nested witness address type",
		masterPriv: "tprv8ZgxMBicQKsPeWwrFuNjEGTTDSY4mRLwd2KDJAPGa1AY" +
			"quw38bZqNMSuB3V1Va3hqJBo9Pt8Sx7kBQer5cNMrb8SYquoWPt9" +
			"Y3BZdhdtUcw",
		accountIndex:       0,
		addrType:           waddrmgr.NestedWitnessPubKey,
		expectedScope:      waddrmgr.KeyScopeBIP0049Plus,
		expectedAddr:       "QdgHt4hN289wZbVWqmHoEPUtPA8joPwMW9",
		expectedChangeAddr: "QNSDM3VJ4Ecohm6PtUwjCabienX7BfKEd5",
	}, {
		name: "bip44 with witness address type",
		masterPriv: "tprv8ZgxMBicQKsPeWwrFuNjEGTTDSY4mRLwd2KDJAPGa1AY" +
			"quw38bZqNMSuB3V1Va3hqJBo9Pt8Sx7kBQer5cNMrb8SYquoWPt9" +
			"Y3BZdhdtUcw",
		accountIndex:       777,
		addrType:           waddrmgr.WitnessPubKey,
		expectedScope:      waddrmgr.KeyScopeBIP0084,
		expectedAddr:       "tltc1qnhgyn2z67ye50ns4huxuvfhghmktzgr20a562q",
		expectedChangeAddr: "tltc1qxyqd8jvr2g6vy8qmj08m6mmr6vm688lz5kvm8p",
	}, {
		name: "traditional bip49",
		masterPriv: "uprv8tXDerPXZ1QsVp8y6GAMSMYxPQgWi3LSY8qS5ZH9x1YRu" +
			"1kGPFjPzR73CFSbVUhdEwJbtsUgucUJ4hGQoJnNepp3RBcE6Jhdom" +
			"FD2KeY6G9",
		accountIndex:       9,
		addrType:           waddrmgr.NestedWitnessPubKey,
		expectedScope:      waddrmgr.KeyScopeBIP0049Plus,
		expectedAddr:       "QhacLWS8hgsD493ZUKuDzyJGC7idz8yZTr",
		expectedChangeAddr: "QdZF5yJGkVJXbgXX4yDdZBiLtTbrcqYP7Y",
	}, {
		name: "bip49+",
		masterPriv: "uprv8tXDerPXZ1QsVp8y6GAMSMYxPQgWi3LSY8qS5ZH9x1YRu" +
			"1kGPFjPzR73CFSbVUhdEwJbtsUgucUJ4hGQoJnNepp3RBcE6Jhdom" +
			"FD2KeY6G9",
		accountIndex:       9,
		addrType:           waddrmgr.WitnessPubKey,
		expectedScope:      waddrmgr.KeyScopeBIP0049Plus,
		expectedAddr:       "QhacLWS8hgsD493ZUKuDzyJGC7idz8yZTr",
		expectedChangeAddr: "tltc1qpys8g8l7duwjvlmzkeyfpcyqxfzsv95fc7zz75",
	}, {
		name: "bip84",
		masterPriv: "vprv9DMUxX4ShgxMM7L5vcwyeSeTZNpxefKwTFMerxB3L1vJ" +
			"x7ZVdutxcUmBDTQBVPMYeaRQeM5FNGpqwysyX1CPT4VeHXJegDX8" +
			"5VJrQvaFaz3",
		accountIndex:       1,
		addrType:           waddrmgr.WitnessPubKey,
		expectedScope:      waddrmgr.KeyScopeBIP0084,
		expectedAddr:       "tltc1qgzww52d2jswxl5zmnltq0d30txtvyvhy8flfqx",
		expectedChangeAddr: "tltc1q6erkczs550rzemk8ut54aekhyla5tuqdk8vc08",
	}}
)

// TestImportAccount tests that extended public keys can successfully be
// imported into both watch only and normal wallets.
func TestImportAccount(t *testing.T) {
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			w, cleanup := testWallet(t)
			defer cleanup()

			testImportAccount(t, w, tc, false, tc.name)
		})

		name := tc.name + " watch-only"
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w, cleanup := testWalletWatchingOnly(t)
			defer cleanup()

			testImportAccount(t, w, tc, true, name)
		})
	}
}

func testImportAccount(t *testing.T, w *Wallet, tc *testCase, watchOnly bool,
	name string) {

	// First derive the master public key of the account we want to import.
	root, err := hdkeychain.NewKeyFromString(tc.masterPriv)
	require.NoError(t, err)

	// Derive the extended private and public key for our target account.
	acct1Pub := deriveAcctPubKey(
		t, root, tc.expectedScope, hardenedKey(tc.accountIndex),
	)

	// We want to make sure we can import and handle multiple accounts, so
	// we create another one.
	acct2Pub := deriveAcctPubKey(
		t, root, tc.expectedScope, hardenedKey(tc.accountIndex+1),
	)

	// And we also want to be able to import loose extended public keys
	// without needing to specify an explicit scope.
	acct3ExternalExtPub := deriveAcctPubKey(
		t, root, tc.expectedScope, hardenedKey(tc.accountIndex+2), 0, 0,
	)
	acct3ExternalPub, err := acct3ExternalExtPub.ECPubKey()
	require.NoError(t, err)

	// Do a dry run import first and check that it results in the expected
	// addresses being derived.
	_, extAddrs, intAddrs, err := w.ImportAccountDryRun(
		name+"1", acct1Pub, root.ParentFingerprint(), &tc.addrType, 1,
	)
	require.NoError(t, err)
	require.Len(t, extAddrs, 1)
	require.Equal(t, tc.expectedAddr, extAddrs[0].Address().String())
	require.Len(t, intAddrs, 1)
	require.Equal(t, tc.expectedChangeAddr, intAddrs[0].Address().String())

	// Import the extended public keys into new accounts.
	acct1, err := w.ImportAccount(
		name+"1", acct1Pub, root.ParentFingerprint(), &tc.addrType,
	)
	require.NoError(t, err)
	require.Equal(t, tc.expectedScope, acct1.KeyScope)

	acct2, err := w.ImportAccount(
		name+"2", acct2Pub, root.ParentFingerprint(), &tc.addrType,
	)
	require.NoError(t, err)
	require.Equal(t, tc.expectedScope, acct2.KeyScope)

	err = w.ImportPublicKey(acct3ExternalPub, tc.addrType)
	require.NoError(t, err)

	// If the wallet is watch only, there is no default account and our
	// imported account will be index 0.
	firstAccountIndex := uint32(1)
	if watchOnly {
		firstAccountIndex = 0
	}

	// We should have 3 additional accounts now.
	acctResult, err := w.Accounts(tc.expectedScope)
	require.NoError(t, err)
	require.Len(t, acctResult.Accounts, int(firstAccountIndex*2)+2)

	// Validate the state of the accounts.
	require.Equal(t, firstAccountIndex, acct1.AccountNumber)
	require.Equal(t, name+"1", acct1.AccountName)
	require.Equal(t, true, acct1.IsWatchOnly)
	require.Equal(t, root.ParentFingerprint(), acct1.MasterKeyFingerprint)
	require.NotNil(t, acct1.AccountPubKey)
	require.Equal(t, acct1Pub.String(), acct1.AccountPubKey.String())
	require.Equal(t, uint32(0), acct1.InternalKeyCount)
	require.Equal(t, uint32(0), acct1.ExternalKeyCount)
	require.Equal(t, uint32(0), acct1.ImportedKeyCount)

	require.Equal(t, firstAccountIndex+1, acct2.AccountNumber)
	require.Equal(t, name+"2", acct2.AccountName)
	require.Equal(t, true, acct2.IsWatchOnly)
	require.Equal(t, root.ParentFingerprint(), acct2.MasterKeyFingerprint)
	require.NotNil(t, acct2.AccountPubKey)
	require.Equal(t, acct2Pub.String(), acct2.AccountPubKey.String())
	require.Equal(t, uint32(0), acct2.InternalKeyCount)
	require.Equal(t, uint32(0), acct2.ExternalKeyCount)
	require.Equal(t, uint32(0), acct2.ImportedKeyCount)

	// Test address derivation.
	extAddr, err := w.NewAddress(acct1.AccountNumber, tc.expectedScope)
	require.NoError(t, err)
	require.Equal(t, tc.expectedAddr, extAddr.String())
	intAddr, err := w.NewChangeAddress(acct1.AccountNumber, tc.expectedScope)
	require.NoError(t, err)
	require.Equal(t, tc.expectedChangeAddr, intAddr.String())

	// Make sure the key count was increased.
	acct1, err = w.AccountProperties(tc.expectedScope, acct1.AccountNumber)
	require.NoError(t, err)
	require.Equal(t, uint32(1), acct1.InternalKeyCount)
	require.Equal(t, uint32(1), acct1.ExternalKeyCount)
	require.Equal(t, uint32(0), acct1.ImportedKeyCount)

	// Make sure we can't get private keys for the imported accounts.
	_, err = w.DumpWIFPrivateKey(intAddr)
	require.True(t, waddrmgr.IsError(err, waddrmgr.ErrWatchingOnly))

	// Get the address info for the single key we imported.
	switch tc.addrType {
	case waddrmgr.NestedWitnessPubKey:
		witnessAddr, err := ltcutil.NewAddressWitnessPubKeyHash(
			ltcutil.Hash160(acct3ExternalPub.SerializeCompressed()),
			&chaincfg.TestNet4Params,
		)
		require.NoError(t, err)

		witnessProg, err := txscript.PayToAddrScript(witnessAddr)
		require.NoError(t, err)

		intAddr, err = ltcutil.NewAddressScriptHash(
			witnessProg, &chaincfg.TestNet4Params,
		)
		require.NoError(t, err)

	case waddrmgr.WitnessPubKey:
		intAddr, err = ltcutil.NewAddressWitnessPubKeyHash(
			ltcutil.Hash160(acct3ExternalPub.SerializeCompressed()),
			&chaincfg.TestNet4Params,
		)
		require.NoError(t, err)

	default:
		t.Fatalf("unhandled address type %v", tc.addrType)
	}

	addrManaged, err := w.AddressInfo(intAddr)
	require.NoError(t, err)
	require.Equal(t, true, addrManaged.Imported())
}
