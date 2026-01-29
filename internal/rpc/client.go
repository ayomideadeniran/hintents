package rpc

import (
	"fmt"

	"github.com/stellar/go/clients/horizonclient"
)

// Client handles interactions with the Stellar Network
type Client struct {
	Horizon *horizonclient.Client
}

// NewClient creates a new RPC client for the specified network
func NewClient(network string) (*Client, error) {
	var horizon *horizonclient.Client

	switch network {
	case "testnet":
		horizon = horizonclient.DefaultTestNetClient
	case "mainnet", "public":
		horizon = horizonclient.DefaultPublicNetClient
	default:
		return nil, fmt.Errorf("unsupported network: %s (use 'testnet' or 'mainnet')", network)
	}

	return &Client{
		Horizon: horizon,
	}, nil
}

// GetTransaction fetches the transaction details and returns envelope and result meta XDR
func (c *Client) GetTransaction(hash string) (envelopeXdr, resultMetaXdr string, err error) {
	tx, err := c.Horizon.TransactionDetail(hash)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch transaction: %w", err)
	}

	return tx.EnvelopeXdr, tx.ResultMetaXdr, nil
}
