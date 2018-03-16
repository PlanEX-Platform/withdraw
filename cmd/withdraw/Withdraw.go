package main

import (
	"eth-withdraw/pkg/config"
	"fmt"
	"github.com/howeyc/gopass"
	"os"
	"eth-withdraw/pkg/accounts"
	"eth-withdraw/pkg/logger"
	"runtime"
	"eth-withdraw/pkg/txupdater"
	"github.com/zhooq/go-ethereum/ethclient"
	"github.com/zhooq/go-ethereum/rpc"
	"eth-withdraw/pkg/listener"
	"math/big"
	"eth-withdraw/pkg/withdraw"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
)

const (
	VERSION = "0.01"
)

func main() {
	setup()

	fmt.Printf("Database password: ")
	pass, _ := gopass.GetPasswd()

	if len(pass) > 1 {
		accounts.KEY = string(pass)
	}

	logger.Log.Printf("Server v%s pid=%d started with processes: %d", VERSION, os.Getpid(), runtime.GOMAXPROCS(runtime.NumCPU()))

	conn, err := ethclient.Dial(config.CFG.BlockchainEndpoint)
	client, err := rpc.Dial(config.CFG.BlockchainEndpoint)

	if err != nil {
		fmt.Println("Failed to connect to the Ethereum client: %v", err)
	}

	go txupdater.StartTxUpdating(client)
	go listener.StartListener(client, conn)

	router := httprouter.New()
	router.POST("/withdraw/", withdraw.MakeWithdraw)

	log.Fatal(http.ListenAndServe(":9011", router))

}

func setup() {
	config.CFG = new(config.Config)

	// Prod env
	//config.CFG.BlockchainEndpoint = "/root/.ethereum/geth.ipc"
	//config.CFG.BlockchainEnandasdpoint = "https://mainnet.infura.io/wRAIg3KbD0yXgE89prjQ"
	config.CFG.BlockchainEndpoint = "ws://128.199.45.106:8546"
	//config.CFG.BlockchainEndpoint = "https://rinkeby.infura.io/wRAIg3KbD0yXgE89prjQ"
	//config.CFG.BlockchainEndpoint = "/root/.local/share/io.parity.ethereum/jsonrpc.ipc"
	//config.CFG.BlockchainEndpoint = "ws://mainnet.dagger.matic.network:1884"
	config.CFG.DBAddr = "localhost:5432"
	config.CFG.DBName = "hirama"
	config.CFG.DBUser = "hirama"
	config.CFG.DBPassword = "hirama"
	config.CFG.RequiredConfirmations = 10
	config.CFG.GasPrice = big.NewInt(5000000000)
	config.CFG.GasLimit = big.NewInt(53000)
}
