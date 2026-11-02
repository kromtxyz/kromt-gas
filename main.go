// KROMT Gas Oracle
// Local client-side gas price & estimation tool using public RPCs only.
// No private keys. No credentials.
//
// https://kromt.xyz  |  https://github.com/kromtxyz

package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var defaultRPCs = map[string]string{
	"robinhood": "https://rpc.mainnet.chain.robinhood.com",
	"ethereum":  "https://eth.llamarpc.com",
	"base":      "https://mainnet.base.org",
	"arbitrum":  "https://arb1.arbitrum.io/rpc",
	"optimism":  "https://mainnet.optimism.io",
}

func main() {
	chain := flag.String("chain", "robinhood", "Chain preset: robinhood|ethereum|base|arbitrum|optimism")
	rpcURL := flag.String("rpc", "", "Custom RPC URL (overrides -chain)")
	to := flag.String("to", "", "Optional: destination address for estimateGas")
	data := flag.String("data", "", "Optional: calldata (hex) for estimateGas")
	value := flag.String("value", "0", "Optional: value in wei for estimateGas")
	flag.Parse()

	url := *rpcURL
	if url == "" {
		var ok bool
		url, ok = defaultRPCs[strings.ToLower(*chain)]
		if !ok {
			fmt.Fprintf(os.Stderr, "Unknown chain preset: %s\n", *chain)
			os.Exit(1)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to RPC: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	fmt.Println("KROMT Gas Oracle")
	fmt.Printf("RPC: %s\n\n", url)

	// Suggest gas price (legacy)
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SuggestGasPrice error: %v\n", err)
	} else {
		fmt.Printf("SuggestGasPrice : %s wei  (~%.4f Gwei)\n", gasPrice.String(), toGwei(gasPrice))
	}

	// EIP-1559 style tip + base fee if available
	tip, err := client.SuggestGasTipCap(ctx)
	if err == nil {
		fmt.Printf("SuggestGasTipCap: %s wei  (~%.4f Gwei)\n", tip.String(), toGwei(tip))
	}

	header, err := client.HeaderByNumber(ctx, nil)
	if err == nil && header.BaseFee != nil {
		fmt.Printf("BaseFee         : %s wei  (~%.4f Gwei)\n", header.BaseFee.String(), toGwei(header.BaseFee))
		if tip != nil {
			maxFee := new(big.Int).Add(header.BaseFee, tip)
			fmt.Printf("Est. maxFeePerGas (base+tip): %s wei (~%.4f Gwei)\n", maxFee.String(), toGwei(maxFee))
		}
	}

	// Optional estimateGas
	if *to != "" {
		toAddr := common.HexToAddress(*to)
		msg := ethereum.CallMsg{
			To: &toAddr,
		}
		if *data != "" {
			msg.Data = common.FromHex(*data)
		}
		if *value != "0" {
			v, ok := new(big.Int).SetString(*value, 10)
			if !ok {
				fmt.Fprintf(os.Stderr, "Invalid value: %s\n", *value)
				os.Exit(1)
			}
			msg.Value = v
		}

		gas, err := client.EstimateGas(ctx, msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nEstimateGas error: %v\n", err)
		} else {
			fmt.Printf("\nEstimateGas     : %d gas units\n", gas)
			if gasPrice != nil {
				cost := new(big.Int).Mul(big.NewInt(int64(gas)), gasPrice)
				fmt.Printf("Est. cost (legacy): %s wei\n", cost.String())
			}
		}
	}

	fmt.Println("\nLocal-only tool. Public RPC only. No private keys.")
	fmt.Println("KROMT · https://kromt.xyz · https://github.com/kromtxyz")
}

func toGwei(wei *big.Int) float64 {
	f := new(big.Float).SetInt(wei)
	gwei := new(big.Float).Quo(f, big.NewFloat(1e9))
	v, _ := gwei.Float64()
	return v
}
