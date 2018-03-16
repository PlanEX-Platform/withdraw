package txupdater

import (
	"time"
	"eth-withdraw/pkg/transactions"
	"eth-withdraw/pkg/logger"
	"eth-withdraw/pkg/config"
	"log"
	"context"
	"github.com/zhooq/go-ethereum/common"
	"github.com/zhooq/go-ethereum/rpc"
	"github.com/zhooq/go-ethereum/core/types"
	"eth-withdraw/pkg/util"
)

type rpcTransaction struct {
	tx *types.Transaction
	txExtraInfo
}

type txExtraInfo struct {
	BlockNumber *string
	BlockHash   common.Hash
	From        common.Address
}

func StartTxUpdating(client *rpc.Client) {

	log.Println("TxUpdater Started")

	func() {
		for true {
			time.Sleep(15 * time.Second)
			checkTs(client)
		}
	}()
}

var transactions = &tx.TransactionSchema{}

func checkTs(client *rpc.Client) {

	transactions.Init()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	txs, err := transactions.Unconfirmed(config.CFG.RequiredConfirmations)

	if err != nil {
		logger.Log.Println("DB get tx error ", err)
	}

	var blockNumber string
	err = client.CallContext(ctx, &blockNumber, "eth_blockNumber")
	if err != nil {
		logger.Log.Println("Block number error ", err)
	}

	currentBlock, _ := utils.ParseInt(blockNumber)

	for _, element := range txs {
		var json *rpcTransaction
		err := client.CallContext(ctx, &json, "eth_getTransactionByHash", element.TxId)
		if err != nil {
			logger.Log.Println("Receipt error ", err)
		}
		minedBlock, _ := utils.ParseInt(*json.txExtraInfo.BlockNumber)
		transactions.UpdateConfirmation(element.TxId, currentBlock-minedBlock)
		logger.Log.Println("Update transaction id: ", element.TxId)
	}

}
