package withdraw

import (
	"github.com/zhooq/go-ethereum/ethclient"
	"context"
	"time"
	"github.com/zhooq/go-ethereum/core/types"
	"github.com/zhooq/go-ethereum/common"
	"math/big"
	"eth-withdraw/config"
	"eth-withdraw/accounts"
	"eth-withdraw/transactions"
	"github.com/zhooq/go-ethereum/crypto"
	"eth-withdraw/ciph"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"net/http"
	"eth-withdraw/logger"
	"eth-withdraw/util"
	"io/ioutil"
	"encoding/json"
	"strings"
)

var acs = &accounts.AccountSchema{}
var transactions = &tx.TransactionSchema{}

type withdrawRequest struct {
	PlanexID   string `json:"id"`
	EthAddress string `json:"to"`
	Amount     string `json:"amount"`
}

func MakeWithdraw(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	buf, _ := ioutil.ReadAll(r.Body)
	body := withdrawRequest{}
	json.Unmarshal(buf, &body)

	body.PlanexID = strings.TrimSpace(body.PlanexID)
	body.EthAddress = strings.TrimSpace(body.EthAddress)
	body.Amount = strings.TrimSpace(body.Amount)

	acs.Init()
	temp, err := acs.ByID(body.PlanexID)

	if err != nil {
		fmt.Fprint(w, string(err.Error()))
	} else {
		addr := common.HexToAddress(body.EthAddress)
		amount, _ := utils.ParseBigInt(body.Amount)
		restx, err := proccedWitdraw(&addr, &amount)
		if err != nil {
			fmt.Fprint(w, err)
		} else if restx != nil {
			fmt.Fprint(w, string(restx.Hash().Hex()))
			logger.Log.Println("Request id: %s", temp.PlanexID)
			logger.Log.Println("Request host: %s Request method: %s", r.Host, r.Method)
		} else {
			fmt.Fprint(w, string("fail"))
		}
	}

	//fmt.Fprint(w, string("ok"))

}

func proccedWitdraw(to *common.Address, amount *big.Int) (*types.Transaction, error) {
	transactions.Init()

	var amountWithFee = big.NewInt(0)
	var fee = big.NewInt(0).Mul(config.CFG.GasLimit, config.CFG.GasPrice)
	amountWithFee.Add(amount, fee)

	acc, err := acs.ByAmountRequired(amountWithFee.String())

	if err != nil {
		logger.Log.Println(err)
	} else {
		privkey, err := ciph.Decrypt(acc.KeyStore, acc.Nonce, accounts.KEY)

		if err != nil {
			logger.Log.Println(err)
			return nil, err
		} else {
			txout, err := sendTx(to, amount, privkey)
			if err != nil {
				logger.Log.Println("Error from: ", err)
				return nil, err
			} else {
				logger.Log.Println("Store to DB: ", txout.Hash().Hex())
				transactions.Create(txout.Hash().Hex(), false, "out", 0, acc.PlanexID)
				accBalance, _ := utils.ParseBigInt(acc.Balance)
				accBalance.Sub(&accBalance, amountWithFee)
				acc.Balance = accBalance.String()
				acs.Update(acc)
				return txout, err
			}
		}
	}

	return nil, err
}

func sendTx(to *common.Address, amount *big.Int, privkey string) (*types.Transaction, error) {
	d := time.Now().Add(1000 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	conn, err := ethclient.Dial(config.CFG.BlockchainEndpoint)

	if err != nil {
		fmt.Println("Failed to connect to the Ethereum client: %v", err)
	}

	key, _ := crypto.HexToECDSA(privkey)
	genAddr := crypto.PubkeyToAddress(key.PublicKey)

	nonce, err := conn.NonceAt(ctx, genAddr, nil)

	if err != nil {
		logger.Log.Println("Cant get nonce", err)
	} else {
		rawtx := types.NewTransaction(nonce, *to, amount, config.CFG.GasLimit, config.CFG.GasPrice, nil)

		signTx, err := types.SignTx(rawtx, types.NewEIP155Signer(big.NewInt(1)), key)
		if err != nil {
			logger.Log.Println("Error from signied tx: ", err)
		}

		err = conn.SendTransaction(ctx, signTx)

		if err != nil {
			logger.Log.Println("Error from sending tx: ", err)
		}
		return signTx, err
	}
	return nil, err
}
