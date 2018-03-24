package main

import (
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
	"eth-withdraw/withdraw"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"eth-withdraw/config"
	"github.com/spf13/viper"
)

const (
	VERSION = "0.01"
)

func init() {
	config.Load()
}

func main() {

	fmt.Printf("Database password: ")
	pass, _ := gopass.GetPasswd()

	if len(pass) > 1 {
		accounts.KEY = string(pass)
	}

	logger.Log.Printf("Server v%s pid=%d started with processes: %d", VERSION, os.Getpid(), runtime.GOMAXPROCS(runtime.NumCPU()))

	conn, err := ethclient.Dial(viper.GetString("BlockchainEndpoint"))
	client, err := rpc.Dial(viper.GetString("BlockchainEndpoint"))

	if err != nil {
		fmt.Println("Failed to connect to the Ethereum client: %v", err)
	}

	go txupdater.StartTxUpdating(client)
	go listener.StartListener(client, conn)

	router := httprouter.New()
	router.POST("/withdraw/", withdraw.MakeWithdraw)

	log.Fatal(http.ListenAndServe(":9011", router))

}
