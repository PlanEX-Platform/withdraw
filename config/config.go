package config

type Config struct {
	BlockchainEndpoint    string
	GasPrice              string
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
