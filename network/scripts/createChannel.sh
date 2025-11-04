#!/bin/bash
set -e

# Run this after startNetwork.sh
BASE_DIR="$(cd "$(dirname "$0")/.." && pwd)"
CHANNEL_ARTIFACTS_DIR="$BASE_DIR/channel-artifacts"

export CORE_PEER_LOCALMSPID=ManufacturerMSP
export CORE_PEER_MSPCONFIGPATH=$BASE_DIR/crypto-config/peerOrganizations/manufacturer.example.com/users/Admin@manufacturer.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export FABRIC_CFG_PATH=$BASE_DIR

ORDERER_HOST=localhost:7050
CHANNEL_NAME=supplychannel
CHANNEL_TX=$CHANNEL_ARTIFACTS_DIR/channel.tx

echo "Creating channel $CHANNEL_NAME..."
peer channel create -o $ORDERER_HOST -c $CHANNEL_NAME -f $CHANNEL_TX --outputBlock $CHANNEL_ARTIFACTS_DIR/$CHANNEL_NAME.block

echo "Joining peer0.manufacturer to $CHANNEL_NAME..."
peer channel join -b $CHANNEL_ARTIFACTS_DIR/$CHANNEL_NAME.block

# Update anchor peers (example)
peer channel update -o $ORDERER_HOST -c $CHANNEL_NAME -f $CHANNEL_ARTIFACTS_DIR/ManufacturerMSPanchors.tx

echo "Channel created and peer joined. Repeat join steps for other orgs if needed."
