#!/bin/bash
export FABRIC_CFG_PATH=~/go/src/github.com/hyperledger/fabric-samples/config
export GOPATH=~/go
export PATH=$PATH:/usr/local/go/bin
echo "127.0.0.1 peer0.org1.example.com" | sudo tee -a /etc/hosts
