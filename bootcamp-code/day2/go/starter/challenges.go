package main

import (
	"fmt"
)

// ==========================================
// CHALLENGE 1: BLOCKCHAIN INFO
// ==========================================
func showBlockchainInfo() error {
	var info struct {
		Chain      string  `json:"chain"`
		Blocks     int     `json:"blocks"`
		Difficulty float64 `json:"difficulty"`
	}

	if err := rpc("getblockchaininfo", nil, "", &info); err != nil {
		return err
	}

	fmt.Println("\n=== Blockchain Info ===")
	fmt.Printf("Chain:      %s\n", info.Chain)
	fmt.Printf("Blocks:     %d\n", info.Blocks)
	fmt.Printf("Difficulty: %v\n", info.Difficulty)
	return nil
}

// ==========================================
// CHALLENGE 2: WALLET BALANCE
// ==========================================
func showWalletBalance(wallet string) error {
	// Wallet may already be loaded -> ignore the error safely
	_ = rpc("loadwallet", []any{wallet}, "", nil)

	var balance float64
	if err := rpc("getbalance", nil, wallet, &balance); err != nil {
		return err
	}

	fmt.Printf("\n=== Wallet: %s ===\n", wallet)
	fmt.Printf("Balance: %v BTC\n", balance)
	return nil
}

// ==========================================
// CHALLENGE 3: LIST TRANSACTIONS
// ==========================================
func listTransactions(wallet string, count int) error {
	_ = rpc("loadwallet", []any{wallet}, "", nil)

	var txs []struct {
		Category      string  `json:"category"`
		Amount        float64 `json:"amount"`
		TxID          string  `json:"txid"`
		Confirmations int     `json:"confirmations"`
	}

	// "*" is used to fetch all categories of transactions
	if err := rpc("listtransactions", []any{"*", count}, wallet, &txs); err != nil {
		return err
	}

	fmt.Printf("\n=== Recent Transactions (%s) ===\n", wallet)
	for _, tx := range txs {
		dir := "OUT"
		switch tx.Category {
		case "receive", "generate", "immature":
			dir = "IN "
		}

		fmt.Printf("%s %+.8f BTC | %d confs\n", dir, tx.Amount, tx.Confirmations)
		fmt.Printf("    TXID: %s\n", tx.TxID)
	}
	return nil
}

// ==========================================
// CHALLENGE 4: DECODE TX
// ==========================================
func decodeTransaction(txid string) error {
	var tx struct {
		Vin []struct {
			Coinbase string `json:"coinbase"`
			TxID     string `json:"txid"`
			Vout     int    `json:"vout"`
		} `json:"vin"`
		Vout []struct {
			Value        float64 `json:"value"`
			ScriptPubKey struct {
				Address string `json:"address"`
			} `json:"scriptPubKey"`
		} `json:"vout"`
	}

	// Passing true asks bitcoind for the verbose decoded JSON transaction object
	if err := rpc("getrawtransaction", []any{txid, true}, "", &tx); err != nil {
		return err
	}

	fmt.Printf("\n=== Decoded Transaction: %s... ===\n", txid[:8])
	fmt.Println("  Inputs (Vin):")
	for _, vin := range tx.Vin {
		if vin.Coinbase != "" {
			fmt.Println("    COINBASE (mining reward)")
		} else {
			fmt.Printf("    From: %s... [vout: %d]\n", vin.TxID[:20], vin.Vout)
		}
	}

	fmt.Println("  Outputs (Vout):")
	for _, vout := range tx.Vout {
		fmt.Printf("    %.8f BTC -> %s\n", vout.Value, vout.ScriptPubKey.Address)
	}
	return nil
}

// ==========================================
// CHALLENGE 5: BLOCK DETAILS
// ==========================================
func showBlock(blockhash string) error {
	// If no hash is provided, fetch the latest (best) block hash
	if blockhash == "" {
		if err := rpc("getbestblockhash", nil, "", &blockhash); err != nil {
			return err
		}
	}

	var block struct {
		Height int      `json:"height"`
		Hash   string   `json:"hash"`
		Time   int64    `json:"time"`
		NTx    int      `json:"nTx"`
		Tx     []string `json:"tx"`
	}

	// Passing 1 tells bitcoind to return detailed block info
	if err := rpc("getblock", []any{blockhash, 1}, "", &block); err != nil {
		return err
	}

	fmt.Println("\n=== Block Details ===")
	fmt.Printf("Height:       #%d\n", block.Height)
	fmt.Printf("Hash:         %s...\n", block.Hash[:32])
	fmt.Printf("Time:         %d\n", block.Time)
	fmt.Printf("Transactions: %d\n", block.NTx)

	// If the block contains transactions, print the first one as an explorer sample
	if len(block.Tx) > 0 {
		fmt.Printf("First TXID:   %s\n", block.Tx[0])
		// Let's hook into Challenge 4 to dynamically decode it!
		_ = decodeTransaction(block.Tx[0])
	}

	return nil
}
