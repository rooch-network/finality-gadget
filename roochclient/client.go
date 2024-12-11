package client

import (
	"context"
	"errors"
	"fmt"
	types "github.com/babylonlabs-io/finality-gadget/roochclient/types"
	"math/big"
)

type RoochClient struct {
	chainID   uint64
	transport RoochTransport
}

// RoochClientOptions configuration options for the RoochClient
type RoochClientOptions struct {
	URL       string
	Transport RoochTransport
}

func NewRoochClient(options RoochClientOptions) *RoochClient {
	var transport RoochTransport
	if options.Transport != nil {
		transport = options.Transport
	} else {
		transportOptions := RoochHTTPTransportOptions{
			URL:     options.URL,
			Headers: nil,
		}
		transport = NewRoochHTTPTransport(transportOptions)
	}

	return &RoochClient{
		transport: transport,
	}
}

func (c *RoochClient) GetRpcApiVersion() (string, error) {
	var resp struct {
		Info struct {
			Version string `json:"version"`
		} `json:"info"`
	}

	err := c.transport.Request("rpc.discover", nil, &resp)
	return resp.Info.Version, err
}

func (c *RoochClient) GetChainId() (uint64, error) {
	if c.chainID != 0 {
		return c.chainID, nil
	}

	var result string
	err := c.transport.Request("rooch_getChainID", nil, &result)
	if err != nil {
		return 0, err
	}

	chainID, err := Str2Uint64(result)
	if err != nil {
		return 0, errors.New("invalid chain ID format")
	}

	c.chainID = chainID
	return chainID, nil
}

func (c *RoochClient) GetBlocks(params *types.GetBlocksParams) (*types.PaginatedBlockViews, error) {
	var result types.PaginatedBlockViews
	err := c.transport.Request("rooch_GetBlocks", []interface{}{
		params.Cursor,
		params.Limit,
		params.DescendingOrder,
	}, &result)
	return &result, err
}

func (c *RoochClient) HeaderByNumber(ctx context.Context, number *big.Int) (*Block, error) {
	// TODO handle negative number
	getBlocksParams := &types.GetBlocksParams{
		//Cursor:          fmt.Sprintf("%d", number),
		Limit:           "1",
		DescendingOrder: true,
	}
	if number.Int64() >= 0 {
		getBlocksParams.Cursor = fmt.Sprintf("%d", number.Int64())
	}

	// batch call
	result, err := c.GetBlocks(getBlocksParams)
	if err != nil {
		return nil, err
	}

	l2blocks, err := BlockViewsToBlocks(result.GetItems())
	if err != nil {
		return nil, err
	}
	if len(l2blocks) <= 0 {
		return nil, fmt.Errorf("the Block of the corresponding height does not exist %v", number)
	}

	return &l2blocks[0], nil
	//block_hash, err := hex.DecodeString(l2blocks[0].BlockHash)
	//if err != nil {
	//	return nil, err
	//}
	//return &client.BlockInfo{
	//	Height: l2blocks[0],
	//	Hash:   block_hash,
	//}, nil
}

// TODO support TransactionByHash RPC
// Instead of Mock implements by getting lastest block
func (c *RoochClient) TransactionByHash(ctx context.Context, txHash string) (*Block, error) {
	getBlocksParams := &types.GetBlocksParams{
		Limit:           "1",
		DescendingOrder: true,
	}

	// batch call
	result, err := c.GetBlocks(getBlocksParams)
	if err != nil {
		return nil, err
	}

	l2blocks, err := BlockViewsToBlocks(result.GetItems())
	if err != nil {
		return nil, err
	}
	if len(l2blocks) <= 0 {
		return nil, fmt.Errorf("Get block fails")
	}

	return &l2blocks[0], nil
}
