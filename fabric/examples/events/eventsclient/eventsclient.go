/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"time"
	"encoding/base64"
	"unicode"
	"strconv"
	"database/sql"

	"github.com/hyperledger/fabric/common/crypto"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/common/localmsp"
	// genesisconfig "github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
	"github.com/hyperledger/fabric/common/tools/protolator"
	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/comm"
	config2 "github.com/hyperledger/fabric/core/config"
	"github.com/hyperledger/fabric/msp"
	common2 "github.com/hyperledger/fabric/peer/common"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/orderer"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/spf13/viper"
	_ "github.com/lib/pq"
)

var (
	channelID        string
	serverAddr       string
	clientKeyPath    string
	clientCertPath   string
	serverRootCAPath string
	stake		 string
	seek             int
	quiet            bool
	filtered         bool
	tlsEnabled       bool
	mTlsEnabled      bool

	configure   = true
	initialized = false
	modify      = false
	oldest      = &orderer.SeekPosition{Type: &orderer.SeekPosition_Oldest{Oldest: &orderer.SeekOldest{}}}
	newest      = &orderer.SeekPosition{Type: &orderer.SeekPosition_Newest{Newest: &orderer.SeekNewest{}}}
	maxStop     = &orderer.SeekPosition{Type: &orderer.SeekPosition_Specified{Specified: &orderer.SeekSpecified{Number: math.MaxUint64}}}
)

const (
	OLDEST = -2
	NEWEST = -1

	ROOT = "configtx"
)

const (
	host = "localhost"
	port = 5432
	user = "agriculture"
	password = "pass"
	dbname = "agricultureDB"
)

var logger = flogging.MustGetLogger("eventsclient")

type blockInfo struct {
	blocknum	uint64		`json:"blocknum"`
	datahash	string		`json:"datahash"`
	prehash		string		`json:"prehash"`
	txid		string		`json:"txid"`
	timestamp	time.Time	`json:"time"`
	info		transactionInfo	`json:"data"`
	modify		modifyLocation	`json:"modify"`
}

type transactionInfo struct {
	chaincodeid	string		`json:"chaincodeID"`
	action		string		`json:"action"`
	key		string		`json:"key"`
	GPS_Location	string		`json:"location"`
	temperature	float64		`json:"temperature"`
	humidity	float64		`json:"humidity"`
}

type modifyLocation struct {
	from		string		`json:"from"`
	to		string		`json:"to"`
}


// deliverClient abstracts common interface
// for deliver and deliverfiltered services
type deliverClient interface {
	Send(*common.Envelope) error
	Recv() (*peer.DeliverResponse, error)
}

// eventsClient client to get connected to the
// events peer delivery system
type eventsClient struct {
	client      deliverClient
	signer      crypto.LocalSigner
	tlsCertHash []byte
}

func (r *eventsClient) seekOldest() error {
	return r.client.Send(r.seekHelper(oldest, maxStop))
}

func (r *eventsClient) seekNewest() error {
	return r.client.Send(r.seekHelper(newest, maxStop))
}

func (r *eventsClient) seekSingle(blockNumber uint64) error {
	specific := &orderer.SeekPosition{Type: &orderer.SeekPosition_Specified{Specified: &orderer.SeekSpecified{Number: blockNumber}}}
	return r.client.Send(r.seekHelper(specific, specific))
}

func (r *eventsClient) seekHelper(start *orderer.SeekPosition, stop *orderer.SeekPosition) *common.Envelope {
	env, err := utils.CreateSignedEnvelopeWithTLSBinding(common.HeaderType_DELIVER_SEEK_INFO, channelID, r.signer, &orderer.SeekInfo{
		Start:    start,
		Stop:     stop,
		Behavior: orderer.SeekInfo_BLOCK_UNTIL_READY,
	}, 0, 0, r.tlsCertHash)
	if err != nil {
		panic(err)
	}
	return env
}

func (r *eventsClient) readEventsStream() {
	db, err := sql.Open("postgres", "user=agriculture password=pass dbname=agricultureDB sslmode=disable")
	checkErr(err)
	defer db.Close()

	checkErr(db.Ping())

	for {
		blockInfo := blockInfo{}
		msg, err := r.client.Recv()
		if err != nil {
			logger.Info("Error receiving:", err)
			return
		}

		switch t := msg.Type.(type) {
		case *peer.DeliverResponse_Status:
			logger.Info("Got status ", t)
			return
		// Print received blocks
		case *peer.DeliverResponse_Block:
			if !quiet {
				logger.Info("Received block: ")
				/*err := protolator.DeepMarshalJSON(os.Stdout, t.Block)
				if err != nil {
					fmt.Printf("  Error pretty printing block: %s", err)
				}*/
				fetchBlockInfo(t.Block.Header, &blockInfo)
				//fmt.Println(base64.StdEncoding.EncodeToString(t.Block.Data.Data))
				//protolator.DeepMarshalJSON(os.Stdout, t.Block.Header)

				for _, test := range t.Block.Data.Data {
					fetchTransactionInfo(test, &blockInfo)
				}

				if configure == true {
					fmt.Println("-----------------------------------------------------------------------------------")
					fmt.Println("| Block Information                                                               |")
					fmt.Println("-----------------------------------------------------------------------------------")
					fmt.Printf("| blocknum    | %-65d |\n", blockInfo.blocknum)
					fmt.Printf("| datahash    | %-65s |\n", blockInfo.datahash)
					fmt.Printf("| prehash     | %-65s |\n", blockInfo.prehash)
					fmt.Printf("| txid        | %-65s |\n", blockInfo.txid)
					fmt.Printf("| timestamp   | %-65s |\n", blockInfo.timestamp)
					fmt.Println("-----------------------------------------------------------------------------------")
					fmt.Println()

					storeBlockInfoToPostgre(db, &blockInfo)
				}

				if modify == true {
					fmt.Println("-----------------------------------------------------------------------------------")
					fmt.Println("| Block Information                                                               |")
					fmt.Println("-----------------------------------------------------------------------------------")
					fmt.Printf("| blocknum    | %-65d |\n", blockInfo.blocknum)
					fmt.Printf("| datahash    | %-65s |\n", blockInfo.datahash)
					fmt.Printf("| prehash     | %-65s |\n", blockInfo.prehash)
					fmt.Printf("| txid        | %-65s |\n", blockInfo.txid)
					fmt.Printf("| timestamp   | %-65s |\n", blockInfo.timestamp)
					fmt.Println("-----------------------------------------------------------------------------------")
					fmt.Println("-----------------------------------------------------------------------------------")
					fmt.Println("| Transaction Information                                                         |")
					fmt.Println("-----------------------------------------------------------------------------------")
					fmt.Printf("| chaincodeid | %-65s |\n", blockInfo.info.chaincodeid)
					fmt.Printf("| action      | %-65s |\n", blockInfo.info.action)
					fmt.Printf("| key         | %-65s |\n", blockInfo.info.key)
					fmt.Printf("| from        | %-65s |\n", blockInfo.modify.from)
					fmt.Printf("| to          | %-65s |\n", blockInfo.modify.to)
					fmt.Println("-----------------------------------------------------------------------------------")
					fmt.Println()

					modify = false
					err := storeBlockInfoToPostgre(db, &blockInfo)

					if err == nil {
						storeModificationInfoToPostgre(db, &blockInfo)
					}

				} else if initialized == true {
					fmt.Println("----------------------------------------------------------------------------------------")
					fmt.Println("| Block Information                                                                    |")
					fmt.Println("----------------------------------------------------------------------------------------")
					fmt.Printf("| blocknum         | %-65d |\n", blockInfo.blocknum)
					fmt.Printf("| datahash         | %-65s |\n", blockInfo.datahash)
					fmt.Printf("| prehash          | %-65s |\n", blockInfo.prehash)
					fmt.Printf("| txid             | %-65s |\n", blockInfo.txid)
					fmt.Printf("| timestamp        | %-65s |\n", blockInfo.timestamp)
					fmt.Println("----------------------------------------------------------------------------------------")
					fmt.Println("----------------------------------------------------------------------------------------")
					fmt.Println("| Transaction Information                                                              |")
					fmt.Println("----------------------------------------------------------------------------------------")
					fmt.Printf("| chaincodeid      | %-65s |\n", blockInfo.info.chaincodeid)
					fmt.Printf("| action           | %-65s |\n", blockInfo.info.action)
					fmt.Printf("| key              | %-65s |\n", blockInfo.info.key)
					fmt.Printf("| location         | %-65s |\n", blockInfo.info.GPS_Location)
					fmt.Printf("| temperature(°C)  | %-65g |\n", blockInfo.info.temperature)
					fmt.Printf("| humidity(%%)      | %-65g |\n", blockInfo.info.humidity)
					fmt.Println("----------------------------------------------------------------------------------------")
					fmt.Println()

					err := storeBlockInfoToPostgre(db, &blockInfo)

					if err == nil && blockInfo.info.action != "init" {
						storeProductInfoToPostgre(db, &blockInfo)
					}
				}
			} else {
				logger.Info("Received block: ", t.Block.Header.Number)
			}

		// Print received filtered blocks
		case *peer.DeliverResponse_FilteredBlock:
			if !quiet {
				logger.Info("Received filtered block: ")
				err := protolator.DeepMarshalJSON(os.Stdout, t.FilteredBlock)
				if err != nil {
					fmt.Printf("  Error pretty printing filtered block: %s", err)
				}

				/*var txActions *peer.FilteredTransactionActions
				filteredTxs := t.FilteredBlock.FilteredTransactions

				for _, tx := range filteredTxs {
					txActions = tx.GetTransactionActions()
				}

				eventExists, eventInfo := FindEvent(txActions)
				if eventExists == true {
					logger.Info("Found chaincode event: ", eventInfo)
				}*/

			} else {
				logger.Info("Received filtered block: ", t.FilteredBlock.Number)
			}
		}
	}
}

func (r *eventsClient) seek(s int) error {
	var err error
	switch seek {
	case OLDEST:
		err = r.seekOldest()
	case NEWEST:
		err = r.seekNewest()
	default:
		err = r.seekSingle(uint64(seek))
	}
	return err
}

// Fetch the block information
func fetchBlockInfo(h *common.BlockHeader, blockInfo *blockInfo) {
	hash := h.GetDataHash()
	prehash := h.GetPreviousHash()

	encodeHash := base64.StdEncoding.EncodeToString(hash)
	encodePreHash := base64.StdEncoding.EncodeToString(prehash)

	blockInfo.blocknum = h.GetNumber()
	blockInfo.datahash = encodeHash
	blockInfo.prehash = encodePreHash

	// fmt.Println(blockInfo.blocknum)
	// fmt.Println(encodeHash)
	// fmt.Println(encodePreHash)
}

// Fetch the transaction information
func fetchTransactionInfo(tdata []byte, blockInfo *blockInfo){
	if tdata == nil {
		fmt.Println("Cannot extract payload from nil transaction")
	}
	if env, err := utils.GetEnvelopeFromBlock(tdata); err != nil {
		fmt.Println("Error getting tx from block(%s)", err)
	} else if env != nil {
		// Get the payload from the envelop
		payload, err := utils.GetPayload(env)
		if err != nil {
			fmt.Println("Couldn not extract payload from envelop, err %s", err)
		}

		chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
		if err != nil {
			fmt.Println("Could not extract channel header from envelop, err %s", err)
		}
		timestamp := time.Unix(chdr.GetTimestamp().Seconds, 0)
		blockInfo.timestamp = timestamp
		blockInfo.txid = chdr.GetTxId()

		if common.HeaderType(chdr.Type) == common.HeaderType_ENDORSER_TRANSACTION {
			// fmt.Println("Flag01")
			tx, err := utils.GetTransaction(payload.Data)
			if err != nil {
				fmt.Println("Error unmarshalling transaction payload for block event: %s", err)
			} else {
				// fmt.Println(tx)
				chaincodeActionPayload, err := utils.GetChaincodeActionPayload(tx.Actions[0].Payload)

				if err == nil {
					// fmt.Println(chaincodeActionPayload)

					ref := string(chaincodeActionPayload.GetChaincodeProposalPayload())
					// fmt.Println(ref)
					f := func(c rune) bool {
						return !unicode.IsLetter(c) && !unicode.IsNumber(c)
					}
					// fmt.Println(f)

					stake := strings.FieldsFunc(ref, f)
					// fmt.Println(stake)

					if initialized == false {
						configure = false

						blockInfo.info.chaincodeid  = stake[6]
						blockInfo.info.action       = stake[9]

						initialized = true
					}

					if stake[4] == "createProduct" {
						blockInfo.info.chaincodeid  = stake[2]
						blockInfo.info.action       = stake[4]
						blockInfo.info.key          = stake[5]

						blockInfo.info.GPS_Location = "(" + stake[6] + "." + stake[7] + ", " + stake[8] + "." + stake[9] + ")"

						temperature_string := stake[10] + "." + stake[11]
						humidity_string     := stake[12] + "." + stake[13]

						temperature_float, err := strconv.ParseFloat(temperature_string, 64)
						checkErr(err)
						humidity_float, err := strconv.ParseFloat(humidity_string, 64)
						checkErr(err)

						blockInfo.info.temperature = temperature_float
						blockInfo.info.humidity    = humidity_float

						/*
						paynum, err := strconv.Atoi(stake[4])
						if err == nil {
							blockInfo.info.paynum = paynum
						}
						*/
					} else if stake[4] == "changeProductLocation" {
						blockInfo.info.chaincodeid = stake[2]
						blockInfo.info.action = stake[4]
						blockInfo.info.key = stake[5]
						blockInfo.modify.to = "(" + stake[6] + "." + stake[7] + ", " + stake[8] + "." + stake[9] + ")"

						modify = true
					}
				}
			}
		}
	}
}

func storeBlockInfoToPostgre(db *sql.DB, blockInfo *blockInfo) error {
	stmt, err := db.Prepare("INSERT INTO blockInfo(blocknum, datahash, prehash, txid, timestamp) VALUES($1, $2, $3, $4, $5);")
	_ = stmt
	checkErr(err)

	res, err := stmt.Exec(blockInfo.blocknum, blockInfo.datahash, blockInfo.prehash, blockInfo.txid, blockInfo.timestamp)
	_ = res
	fmt.Println(err)
	fmt.Println()

	// fmt.Println(stmt)
	// fmt.Println(res)

	return err
}

func storeProductInfoToPostgre(db *sql.DB, blockInfo *blockInfo) {
	stmt, err := db.Prepare("INSERT INTO productInfo(blocknum, key, location, temperature, humidity) VALUES($1, $2, $3, $4, $5);")
	_ = stmt
	checkErr(err)

	res, err := stmt.Exec(blockInfo.blocknum, blockInfo.info.key, blockInfo.info.GPS_Location, blockInfo.info.temperature, blockInfo.info.humidity)
	_ = res
	fmt.Println(err)
	fmt.Println()
}

func storeModificationInfoToPostgre(db *sql.DB, blockInfo *blockInfo) {
	stmt, err := db.Prepare("INSERT INTO modifyLocationInfo(blocknum, key, original, modification) VALUES($1, $2, $3, $4);")
	_ = stmt
	checkErr(err)

	res, err := stmt.Exec(blockInfo.blocknum, blockInfo.info.key, blockInfo.modify.from, blockInfo.modify.to)
	_ = res
	fmt.Println(err)
	fmt.Println()
}

func checkErr(err error) {
	if err != nil {
		// fmt.Println(err)
		panic(err)
	}
}

func main() {
	initConfig()
	initMSP()
	readCLInputs()

	if seek < OLDEST {
		logger.Info("Invalid seek value")
		flag.PrintDefaults()
		return
	}

	clientConfig := comm.ClientConfig{
		KaOpts:  comm.DefaultKeepaliveOptions,
		SecOpts: &comm.SecureOptions{},
		Timeout: 5 * time.Minute,
	}

	if tlsEnabled {
		clientConfig.SecOpts.UseTLS = true
		rootCert, err := ioutil.ReadFile(serverRootCAPath)
		if err != nil {
			logger.Info("error loading TLS root certificate", err)
			return
		}
		clientConfig.SecOpts.ServerRootCAs = [][]byte{rootCert}
		if mTlsEnabled {
			clientConfig.SecOpts.RequireClientCert = true
			clientKey, err := ioutil.ReadFile(clientKeyPath)
			if err != nil {
				logger.Info("error loading client TLS key from", clientKeyPath)
				return
			}
			clientConfig.SecOpts.Key = clientKey

			clientCert, err := ioutil.ReadFile(clientCertPath)
			if err != nil {
				logger.Info("error loading client TLS certificate from path", clientCertPath)
				return
			}
			clientConfig.SecOpts.Certificate = clientCert
		}
	}

	grpcClient, err := comm.NewGRPCClient(clientConfig)
	if err != nil {
		logger.Info("Error creating grpc client:", err)
		return
	}
	conn, err := grpcClient.NewConnection(serverAddr, "")
	if err != nil {
		logger.Info("Error connecting:", err)
		return
	}

	var client deliverClient
	if filtered {
		client, err = peer.NewDeliverClient(conn).DeliverFiltered(context.Background())
	} else {
		client, err = peer.NewDeliverClient(conn).Deliver(context.Background())
	}

	if err != nil {
		logger.Info("Error connecting:", err)
		return
	}

	events := &eventsClient{
		client: client,
		signer: localmsp.NewSigner(),
	}

	if mTlsEnabled {
		events.tlsCertHash = util.ComputeSHA256(grpcClient.Certificate().Certificate[0])
	}

	events.seek(seek)
	if err != nil {
		logger.Info("Received error:", err)
		return
	}

	events.readEventsStream()
}

func readCLInputs() {
	flag.StringVar(&serverAddr, "server", "peer0.org1.example.com:7051", "The RPC server to connect to.")
	flag.StringVar(&channelID, "channelID", "mychannel", "The channel ID to deliver from.")
	flag.BoolVar(&quiet, "quiet", false, "Only print the block number, will not attempt to print its block contents.")
	flag.BoolVar(&filtered, "filtered", false, "Whenever to read filtered events from the peer delivery service or get regular blocks.")
	flag.BoolVar(&tlsEnabled, "tls", true, "TLS enabled/disabled")
	flag.BoolVar(&mTlsEnabled, "mTls", true, "Mutual TLS enabled/disabled (whenever server side validates clients TLS certificate)")
	flag.StringVar(&clientKeyPath, "clientKey", "/home/chris/go/src/github.com/hyperledger/fabric-samples/first-network/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.key", "Specify path to the client TLS key")
	flag.StringVar(&clientCertPath, "clientCert", "/home/chris/go/src/github.com/hyperledger/fabric-samples/first-network/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.crt", "Specify path to the client TLS certificate")
	flag.StringVar(&serverRootCAPath, "rootCert", "/home/chris/go/src/github.com/hyperledger/fabric-samples/first-network/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt", "Specify path to the server root CA certificate")
	flag.IntVar(&seek, "seek", OLDEST, "Specify the range of requested blocks."+
		"Acceptable values:"+
		"-2 (or -1) to start from oldest (or newest) and keep at it indefinitely."+
		"N >= 0 to fetch block N only.")
	flag.Parse()
}

/*
func FindEvent(s interface{})(bool, string){
	value := fmt.Sprintf("%v", s)
	if value != "" {
		return true, value
	}
	return false, ""
}*/

func initMSP() {
	// Init the MSP
	var mspMgrConfigDir = config2.GetPath("peer.mspConfigPath")
	var mspID = viper.GetString("peer.localMspId")
	var mspType = viper.GetString("peer.localMspType")

	mspMgrConfigDir = "/home/chris/go/src/github.com/hyperledger/fabric-samples/first-network/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp"
	mspID = "Org1MSP"
	mspType = "bccsp"

	if mspType == "" {
		mspType = msp.ProviderTypeToString(msp.FABRIC)
	}
	err := common2.InitCrypto(mspMgrConfigDir, mspID, mspType)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Sprintf("Cannot run client because %s", err.Error()))
	}
}

func initConfig() {
	// For environment variables.
	viper.SetEnvPrefix(ROOT)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	err := common2.InitConfig(ROOT)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error when initializing %s config : %s", ROOT, err))
	}
}
