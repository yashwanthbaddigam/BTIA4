/*
 A simple Node.js script to submit transactions to the shipment chaincode.
 Assumes you have a wallet with 'manufacturer' identity and connection profile at ../network/connection-org1.json
*/

const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const fs = require('fs');

async function main() {
  try {
    const ccpPath = path.resolve(__dirname, '..', 'network', 'connection-org1.json');
    const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
    const ccp = JSON.parse(ccpJSON);

    const walletPath = path.join(__dirname, 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);

    const identityLabel = 'manufacturer'; // identity must exist in wallet
    const identity = await wallet.get(identityLabel);
    if (!identity) {
      console.log(`Identity ${identityLabel} not found in wallet. Please enroll/import it before running this script.`);
      return;
    }

    const gateway = new Gateway();
    await gateway.connect(ccp, {
      wallet,
      identity: identityLabel,
      discovery: { enabled: true, asLocalhost: true },
    });

    const network = await gateway.getNetwork('supplychannel'); // channel name
    const contract = network.getContract('shipment'); // chaincode name

    console.log('Submitting CreateShipment transaction...');
    await contract.submitTransaction('CreateShipment', 'SHIP1', 'Laptop Model X', 'Manufacturer', new Date().toISOString());
    console.log('CreateShipment committed.');

    console.log('Submitting TransferShipment transaction (to Distributor)...');
    await contract.submitTransaction('TransferShipment', 'SHIP1', 'Distributor', new Date().toISOString());
    console.log('TransferShipment committed.');

    console.log('Submitting ReceiveShipment transaction (retailer receives later)...');
    await contract.submitTransaction('ReceiveShipment', 'SHIP1', new Date().toISOString());
    console.log('ReceiveShipment committed.');

    console.log('Evaluating QueryShipment...');
    const result = await contract.evaluateTransaction('QueryShipment', 'SHIP1');
    console.log('QueryShipment result:', result.toString());

    await gateway.disconnect();
  } catch (error) {
    console.error('Error:', error);
    process.exit(1);
  }
}

main();
