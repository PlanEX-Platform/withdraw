package listener

import (
	"math/big"
	"github.com/zhooq/go-ethereum/ethclient"
	"github.com/zhooq/go-ethereum/common"
	"context"
	"github.com/zhooq/go-ethereum/core/types"
	"time"
	"withdraw/accounts"
	"log"
	"github.com/zhooq/go-ethereum/rpc"
	"withdraw/transactions"
	"withdraw/logger"
	"withdraw/util"
)

type Transaction struct {
	Hash     string
	ValueWei *big.Int
	To       string
}

type BlockP struct {
	Number string
}

type TransactionHandler func(transaction Transaction) error

var acs = &accounts.AccountSchema{}
var transactions = &tx.TransactionSchema{}

func UpdateBalance(tx *Transaction) {
	acs.Init()
	transactions.Init()

	acc, err := acs.ByAddress(tx.To)
	_, err = transactions.ByTxID(tx.Hash)

	// check account in db
	if acc.EthAddress != "" && err != nil {
		logger.Log.Println("Increase balance for: %s", acc.EthAddress)
		amount, _ := utils.ParseBigInt(acc.Balance)
		newAmount, _ := utils.ParseBigInt(tx.ValueWei.String())

		sum := new(big.Int).Add(&amount, &newAmount)
		acc.Balance = sum.String()
		acs.Update(acc)
		transactions.Create(tx.Hash, true, "in", 1, acc.PlanexID)
		//api.MakeTransferToViu(acc.ViulyID, tx, ev)
	}

}

func GetBalance(client *ethclient.Client, account common.Address) (*big.Int, error) {
	res, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		logger.Log.Println("Cannot read balance", err)
	}
	return res, err
}

func StartListener(l *rpc.Client, conn *ethclient.Client) {

	log.Println("Started listener")

	clientSubscription(l, conn)

}

func getBlock(l *ethclient.Client, number *big.Int) {
	d := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	block, _ := l.BlockByNumber(ctx, number)

	processBlock(block)
}

func processBlock(block *types.Block) error {
	transactions := block.Transactions()

	logger.Log.Println("Proceed block: ", block.Number())

	for _, transaction := range transactions {
		to := transaction.To()
		if to == nil {
			// Contract creation
			continue
		}

		txs := Transaction{
			Hash:     transaction.Hash().Hex(),
			ValueWei: transaction.Value(),
			To:       to.Hex(),
		}

		UpdateBalance(&txs)
	}

	return nil
}

func clientSubscription(client *rpc.Client, conn *ethclient.Client) {

	subch := make(chan BlockP)

	// Ensure that subch receives the latest block.
	go func() {
		for i := 0; ; i++ {
			if i > 0 {
				time.Sleep(2 * time.Second)
			}
			subscribeBlocks(client, subch)
		}
	}()

	// Print events from the subscription as they arrive.
	for block := range subch {
		blockNumber, _ := utils.ParseBigInt(block.Number)

		getBlock(conn, &blockNumber)
	}
}

// subscribeBlocks runs in its own goroutine and maintains
// a subscription for new blocks.
func subscribeBlocks(client *rpc.Client, subch chan BlockP) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Subscribe to new blocks.
	sub, err := client.EthSubscribe(ctx, subch, "newHeads")
	if err != nil {
		logger.Log.Println("subscribe error:", err)
		return
	}

	// The connection is established now.
	// Update the channel with the current block.
	var lastBlock BlockP
	if err := client.CallContext(ctx, &lastBlock, "eth_getBlockByNumber", "latest", true); err != nil {
		logger.Log.Println("can't get latest block:", err)
		return
	}
	subch <- lastBlock

	// The subscription will deliver events to the channel. Wait for the
	// subscription to end for any reason, then loop around to re-establish
	// the connection.
	logger.Log.Println("connection lost: ", <-sub.Err())
}
