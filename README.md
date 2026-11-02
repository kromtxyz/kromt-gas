# KROMT Gas Oracle

**Local gas price and estimation utility for the KROMT ecosystem.**

Part of the public tooling suite for [KROMT](https://kromt.xyz) — the non-custodial AI trading agent.

- Website: https://kromt.xyz  
- GitHub: https://github.com/kromtxyz  

---

## What it does

Fetches current gas market data and optionally estimates gas for a transaction using **public RPC endpoints only**.

- SuggestGasPrice (legacy)
- SuggestGasTipCap (EIP-1559)
- Current base fee
- Optional `eth_estimateGas` for a given `to` + calldata

No private keys are ever used or required.

## Requirements

- Go 1.21+
- `github.com/ethereum/go-ethereum`

```bash
go mod tidy
go build -o kromt-gas .
```

## Usage

```bash
# Default – Robinhood Chain
./kromt-gas

# Other chains
./kromt-gas -chain ethereum
./kromt-gas -chain base

# Custom RPC
./kromt-gas -rpc https://your.public.rpc

# Estimate gas for a call
./kromt-gas -to 0xContractAddress -data 0xa9059cbb...
./kromt-gas -to 0xContract -data 0x... -value 1000000000000000000
```

## Security Notes

- Read-only tool.
- Connects only to the RPC URL you specify (or the public defaults).
- Never handles private keys, mnemonics, or signing.

## License

MIT — see [LICENSE](LICENSE)

---

**KROMT** · Non-custodial AI trading agent  
https://kromt.xyz · https://github.com/kromtxyz
