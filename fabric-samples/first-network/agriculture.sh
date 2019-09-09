sudo ./byfn.sh down ;
sudo ../bin/cryptogen generate --config=./crypto-config.yaml ;
export FABRIC_CFG_PATH=$PWD ;
sudo ../bin/configtxgen -profile TwoOrgsOrdererGenesis -channelID byfn-sys-channel -outputBlock ./channel-artifacts/genesis.block ;
export CHANNEL_NAME=mychannel && sudo ../bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME ;
sudo ../bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP ;
sudo ../bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP ;
sudo docker-compose -f docker-compose-cli.yaml up -d ;
sudo docker exec -it cli bash ;
