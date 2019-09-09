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

