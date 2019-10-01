#!/usr/bin/env bash
while IFS=',' read -r key lng lat temperature humidity; do
    location_str="("$lng", "$lat")"
    echo "$key"
    echo "$location_str"
    echo "$temperature"
    echo "$humidity"

    arg_str="{"\""Args"\"":["\""createProduct"\"", "\""$key"\"", "\""$location_str"\"", "\""$temperature"\"", "\""$humidity"\""]}"

    peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n mycc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c "$arg_str"
    sleep 2
done < file.csv