package main

import (
	"eth-withdraw/config"
	"fmt"
	"github.com/howeyc/gopass"
	"os"
	"eth-withdraw/accounts"
	"eth-withdraw/logger"
	"runtime"
	"eth-withdraw/txupdater"
	"github.com/zhooq/go-ethereum/ethclient"
	"github.com/zhooq/go-ethereum/rpc"
	"eth-withdraw/listener"
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
	listener.StartListener(client, conn)

}

func setup() {
	config.CFG = new(config.Config)

	// Prod env
	//config.CFG.BlockchainEndpoint = "/root/.ethereum/geth.ipc"
	//config.CFG.BlockchainEndpoint = "https://mainnet.infura.io/wRAIg3KbD0yXgE89prjQ"
	config.CFG.BlockchainEndpoint = "ws://128.199.45.106:8546"
	//config.CFG.BlockchainEndpoint = "https://rinkeby.infura.io/wRAIg3KbD0yXgE89prjQ"
	//config.CFG.BlockchainEndpoint = "/root/.local/share/io.parity.ethereum/jsonrpc.ipc"
	//config.CFG.BlockchainEndpoint = "ws://mainnet.dagger.matic.network:1884"
	config.CFG.GasPrice = "10000000000"
	config.CFG.DBAddr = "localhost:5432"
	config.CFG.DBName = "hirama"
	config.CFG.DBUser = "hirama"
	config.CFG.DBPassword = "hirama"
	config.CFG.RequiredConfirmations = 10
}
