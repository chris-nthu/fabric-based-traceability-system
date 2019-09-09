#!/bin/bash
#FABRIC_CFG_PATH=/home/chris/go/src/github.com/hyperledger/fabric-samples/config \
#CORE_PEER_LOCALMSPID="Org1MSP" \
#CORE_PEER_MSPCONFIGPATH=/home/chris/go/src/github.com/hyperledger/fabric-samples/first-network/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp \
#./eventsclient -server=127.0.0.1:7051 -channelID=mychannel -filtered=true -tls=true \
# -clientKey=/home/chris/go/src/github.com/hyperledger/fabric-samples/first-network/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.key \
# -clientCert=/home/chris/go/src/github.com/hyperledger/fabric-samples/first-network/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.crt \
# -rootCert=/home/chris/go/src/github.com/hyperledger/fabric-samples/first-network/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/ca.crt

go build ;
sudo FABRIC_CFG_PATH=/home/chris/go/src/github.com/hyperledger/fabric-samples/config \
./eventsclient
