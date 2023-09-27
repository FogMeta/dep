#!/bin/bash
export AKASH_KEY_NAME=$1
echo "account: $AKASH_KEY_NAME"
export AKASH_KEYRING_BACKEND=os
export AKASH_ACCOUNT_ADDRESS="$(provider-services keys show $AKASH_KEY_NAME -a)"
echo "wallet: $AKASH_ACCOUNT_ADDRESS"
export AKASH_NET="https://raw.githubusercontent.com/akash-network/net/main/mainnet"
export AKASH_VERSION="$(curl -s https://api.github.com/repos/akash-network/provider/releases/latest | jq -r '.tag_name')"
export AKASH_CHAIN_ID="$(curl -s "$AKASH_NET/chain-id.txt")"
export AKASH_NODE="$(curl -s "$AKASH_NET/rpc-nodes.txt" | shuf -n 1)"
echo "node: $AKASH_NODE, chain id: $AKASH_CHAIN_ID, backend: $AKASH_KEYRING_BACKEND"
export AKASH_GAS=auto
export AKASH_GAS_ADJUSTMENT=1.25
export AKASH_GAS_PRICES=0.025uakt
export AKASH_SIGN_MODE=amino-json