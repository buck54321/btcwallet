package chain

import (
	"container/list"
	"errors"

	"github.com/dcrlabs/ltcwallet/spv"
	"github.com/dcrlabs/ltcwallet/spv/banman"
	"github.com/dcrlabs/ltcwallet/spv/headerfs"
	"github.com/ltcsuite/ltcd/chaincfg"
	"github.com/ltcsuite/ltcd/chaincfg/chainhash"
	"github.com/ltcsuite/ltcd/ltcutil"
	"github.com/ltcsuite/ltcd/ltcutil/gcs"
	"github.com/ltcsuite/ltcd/wire"
	"github.com/stretchr/testify/mock"
)

var (
	errNotImplemented = errors.New("not implemented")
	testBestBlock     = &headerfs.BlockStamp{
		Height: 42,
	}
)

var (
	_ rescanner            = (*mockRescanner)(nil)
	_ NeutrinoChainService = (*mockChainService)(nil)
)

// newMockNeutrinoClient constructs a neutrino client with a mock chain
// service implementation and mock rescanner interface implementation.
func newMockNeutrinoClient() *NeutrinoClient {
	// newRescanFunc returns a mockRescanner
	newRescanFunc := func(ro ...spv.RescanOption) rescanner {
		return &mockRescanner{
			updateArgs: list.New(),
		}
	}

	return &NeutrinoClient{
		CS:        &mockChainService{},
		newRescan: newRescanFunc,
	}
}

// mockRescanner is a mock implementation of a rescanner interface for use in
// tests.  Only the Update method is implemented.
type mockRescanner struct {
	updateArgs *list.List
}

func (m *mockRescanner) Update(opts ...spv.UpdateOption) error {
	m.updateArgs.PushBack(opts)
	return nil
}

func (m *mockRescanner) Start() <-chan error {
	return nil
}

func (m *mockRescanner) WaitForShutdown() {
	// no-op
}

// mockChainService is a mock implementation of a chain service for use in
// tests.  Only the Start, GetBlockHeader and BestBlock methods are implemented.
type mockChainService struct {
}

func (m *mockChainService) Start() error {
	return nil
}

func (m *mockChainService) BestBlock() (*headerfs.BlockStamp, error) {
	return testBestBlock, nil
}

func (m *mockChainService) GetBlockHeader(
	*chainhash.Hash) (*wire.BlockHeader, error) {

	return &wire.BlockHeader{}, nil
}

func (m *mockChainService) GetBlock(chainhash.Hash,
	...spv.QueryOption) (*ltcutil.Block, error) {

	return nil, errNotImplemented
}

func (m *mockChainService) GetBlockHeight(*chainhash.Hash) (int32, error) {
	return 0, errNotImplemented
}

func (m *mockChainService) GetBlockHash(int64) (*chainhash.Hash, error) {
	return nil, errNotImplemented
}

func (m *mockChainService) IsCurrent() bool {
	return false
}

func (m *mockChainService) SendTransaction(*wire.MsgTx) error {
	return errNotImplemented
}

func (m *mockChainService) GetCFilter(chainhash.Hash,
	wire.FilterType, ...spv.QueryOption) (*gcs.Filter, error) {

	return nil, errNotImplemented
}

func (m *mockChainService) GetUtxo(
	_ ...spv.RescanOption) (*spv.SpendReport, error) {

	return nil, errNotImplemented
}

func (m *mockChainService) BanPeer(string, banman.Reason) error {
	return errNotImplemented
}

func (m *mockChainService) IsBanned(addr string) bool {
	panic(errNotImplemented)
}

func (m *mockChainService) AddPeer(*spv.ServerPeer) {
	panic(errNotImplemented)
}

func (m *mockChainService) AddBytesSent(uint64) {
	panic(errNotImplemented)
}

func (m *mockChainService) AddBytesReceived(uint64) {
	panic(errNotImplemented)
}

func (m *mockChainService) NetTotals() (uint64, uint64) {
	panic(errNotImplemented)
}

func (m *mockChainService) UpdatePeerHeights(*chainhash.Hash,
	int32, *spv.ServerPeer,
) {
	panic(errNotImplemented)
}

func (m *mockChainService) ChainParams() chaincfg.Params {
	panic(errNotImplemented)
}

func (m *mockChainService) Stop() error {
	panic(errNotImplemented)
}

func (m *mockChainService) PeerByAddr(string) *spv.ServerPeer {
	panic(errNotImplemented)
}

// mockRPCClient mocks the rpcClient interface.
type mockRPCClient struct {
	mock.Mock
}

// Compile time assertion that MockPeer implements lnpeer.Peer.
var _ rpcClient = (*mockRPCClient)(nil)

func (m *mockRPCClient) GetRawMempool() ([]*chainhash.Hash, error) {
	args := m.Called()
	return args.Get(0).([]*chainhash.Hash), args.Error(1)
}

func (m *mockRPCClient) GetRawTransaction(
	txHash *chainhash.Hash) (*ltcutil.Tx, error) {

	args := m.Called(txHash)
	return args.Get(0).(*ltcutil.Tx), args.Error(1)
}
