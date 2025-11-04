# Supply Chain Blockchain Network (ICT4415)

## Overview
A simple Hyperledger Fabric network simulating shipment tracking among Manufacturer, Distributor, and Retailer.

## Setup
1. Navigate to `/network` and run:
   ```bash
   ./scripts/startNetwork.sh
   ./scripts/createChannel.sh
   ```
2. Deploy chaincode:
   ```bash
   peer lifecycle chaincode install shipment.tar.gz
   ```
3. Run the simulation:
   ```bash
   node application/app.js
   ```

## Security Analysis
See `/report/Security_Threat_Analysis.pdf` for risks, mitigations, and an attack example.
