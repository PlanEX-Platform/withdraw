package accounts

import (
	"encoding/hex"
	"fmt"
	"eth-withdraw/ciph"
	"eth-withdraw/config"
	"github.com/go-pg/pg"
	"github.com/rkuris/go.uuid"
	"github.com/zhooq/go-ethereum/crypto"
)

var KEY string

type Account struct {
	ID         string
	PlanexID    string
	EthAddress string
	Balance    string
	KeyStore   string
	Nonce      string
}

func (a Account) String() string {
	return fmt.Sprintf("Account <%s %s %s %s>", a.ID, a.PlanexID, a.EthAddress, a.Balance)
}

/**
Generate new eth account
*/
func GetAccount() (string, string) {

	p, _ := crypto.GenerateKey()

	private := hex.EncodeToString(crypto.FromECDSA(p))
	public := crypto.PubkeyToAddress(p.PublicKey).Hex()

	return public, private
}

type AccountSchema struct {
	db *pg.DB
}

func (schema *AccountSchema) Init() (*pg.DB, error) {
	var db *pg.DB
	if schema.db == nil {
		db := pg.Connect(&pg.Options{
			Addr:     config.CFG.DBAddr,
			Database: config.CFG.DBName,
			User:     config.CFG.DBUser,
			Password: config.CFG.DBPassword,
		})
		schema.db = db
		for _, model := range []interface{}{&Account{}} {
			err := db.CreateTable(model, nil)
			if nil != err {
				return nil, err
			}
		}
	}
	return db, nil
}

func (schema *AccountSchema) Create(planexID string, ethAddress string, PrivKey string) (Account, error) {
	ciphText, nonce, _ := ciph.Encrypt(PrivKey, KEY)
	newAcc := &Account{
		ID:         uuid.NewV4().String(),
		PlanexID:    planexID,
		EthAddress: ethAddress,
		Balance:    "0",
		KeyStore:   ciphText,
		Nonce:      nonce}
	err := schema.db.Insert(newAcc)
	return *newAcc, err
}

func (schema AccountSchema) All() ([]Account, error) {
	var accounts []Account
	err := schema.db.Model(&accounts).Select()
	return accounts, err
}

func (schema *AccountSchema) ByID(planexID string) (Account, error) {
	var acc = &Account{}
	err := schema.db.Model(acc).
		Where("planex_id = ?", planexID).
		Select()
	if acc.EthAddress == "" {
		newAddr, PrivKey := GetAccount()
		newAcc, _ := schema.Create(planexID, newAddr, PrivKey)
		acc = &newAcc
	}
	return *acc, err
}

func (schema *AccountSchema) ByAddress(ethAddress string) (Account, error) {
	account := Account{}
	err := schema.db.Model(&account).
		Where("eth_address = ?", ethAddress).
		Select()
	return account, err
}

func (schema AccountSchema) Update(acc Account) error {
	err := schema.db.Update(&acc)
	return err
}

func (schema AccountSchema) Delete(acc Account) error {
	err := schema.db.Delete(&acc)
	return err
}
