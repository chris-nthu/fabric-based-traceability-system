#!/bin/sh

sudo docker exec cli peer chaincode query -C mychannel -n mycc -c '{"Args":["queryProduct", "No998"]}' ;