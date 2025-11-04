#!/bin/bash
set -e

# Paths (assumes script run from network/scripts)
BASE_DIR="$(cd "$(dirname "$0")/.." && pwd)"
CRYPTO_CONFIG="$BASE_DIR/crypto-config.yaml"
CHANNEL_ARTIFACTS_DIR="$BASE_DIR/channel-artifacts"

echo "Cleaning previous artifacts..."
rm -rf $BASE_DIR/crypto-config $BASE_DIR/crypto-config-orderer $CHANNEL_ARTIFACTS_DIR
mkdir -p $CHANNEL_ARTIFACTS_DIR

echo "Generating crypto material with cryptogen..."
which cryptogen >/dev/null 2>&1 || { echo "cryptogen not found in PATH"; exit 1; }
cryptogen generate --config=$CRYPTO_CONFIG --output="$BASE_DIR/crypto-config"

echo "Generating genesis block and channel config (configtxgen)..."
which configtxgen >/dev/null 2>&1 || { echo "configtxgen not found in PATH"; exit 1; }

export FABRIC_CFG_PATH=$BASE_DIR
configtxgen -profile SampleSingleChannel -channelID system-channel -outputBlock $CHANNEL_ARTIFACTS_DIR/genesis.block
configtxgen -profile SampleSingleChannel -outputCreateChannelTx $CHANNEL_ARTIFACTS_DIR/channel.tx -channelID supplychannel
configtxgen -profile SampleSingleChannel -outputAnchorPeersUpdate $CHANNEL_ARTIFACTS_DIR/ManufacturerMSPanchors.tx -channelID supplychannel -asOrg ManufacturerMSP

echo "Bringing up network containers (docker-compose)..."
cd $BASE_DIR
docker-compose -f docker-compose.yaml up -d

echo "Network started. You may now create the channel using createChannel.sh"
