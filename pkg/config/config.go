package config

import "math/big"

type Config struct {
	BlockchainEndpoint    string
	GasPrice              *big.Int
	GasLimit              *big.Int
	DBAddr                string
	DBName                string
	DBUser                string
	DBPassword            string
	EmailServer           string
	EmailPort             int
	EmailUser             string
	EmailPassword         string
	RequiredConfirmations uint
}

var CFG *Config
