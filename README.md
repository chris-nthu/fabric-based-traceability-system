# fabric-based-traceability-system

An agricultural tracebility system made by hyperledger fabric.

## BACKUP

Use `collect_files.sh` to move the files I want to backup to this directory.

```bash
$ sudo ./collect_files.sh
```

## PUSH TO GITHUB

```bash
$ git add -A
$ git commit -m "Update message"
$ git push
```

## USAGE

1. Create our own blockchain system

```bash
/home/user/go/src/github.com/hyperledger/fabric-samples/first-network$ sudo ./agriculture.sh
```

> After that, you will enter the cli container.

2. Set some configuration and install chaincode

```bash
opt/gopath/src/github.com/hyperledger/fabric/peer# ./../../../chaincode/agriculture/settings.sh
```

3. Set channel name

```bash
opt/gopath/src/github.com/hyperledger/fabric/peer# export CHANNEL_NAME=mychannel
```

4. (Option) Test submit lots of transaction

```bash
opt/gopath/src/github.com/hyperledger/fabric/peer# ./../../../chaincode/agriculture/test/test.sh
```

5. Execute eventcilent to fetch transaction and store into PostgreSQL

```bash
/home/user/go/src/github.com/hyperledger/fabric/examples/events/eventsclient$ go build

/home/user/go/src/github.com/hyperledger/fabric/examples/events/eventsclient$ sudo FABRIC_CFG_PATH=$GOPATH/src/github.com/hyperledger/fabric-samples/config ./eventsclient
```
