package tx

import (
	"fmt"
	"eth-withdraw/config"
	"github.com/go-pg/pg"
)

type Transaction struct {
	ID            int
	TxId          string
	Mined         bool
	TxType        string
	Confirmations int
	AccountID     string
	Confirmed     bool
}

func (t Transaction) String() string {
	return fmt.Sprintf("Transaction <%s %b %s %d %s %b>",
		t.TxId,
		t.Mined,
		t.TxType,
		t.Confirmations,
		t.AccountID,
		t.Confirmed)
}

type TransactionSchema struct {
	db *pg.DB
}

func (schema *TransactionSchema) Init() (*pg.DB, error) {
	var db *pg.DB
	if schema.db == nil {
		db := pg.Connect(&pg.Options{
			Addr:     config.CFG.DBAddr,
			Database: config.CFG.DBName,
			User:     config.CFG.DBUser,
			Password: config.CFG.DBPassword,
		})
		schema.db = db
		for _, model := range []interface{}{&Transaction{}} {
			err := db.CreateTable(model, nil)
			if nil != err {
				return nil, err
			}
		}
	}
	return db, nil
}

func (schema *TransactionSchema) Create(tx string, mined bool, txType string, confirmation int, accountId string) (Transaction, error) {
	newTX := &Transaction{
		TxId:          tx,
		Mined:         mined,
		TxType:        txType,
		Confirmations: confirmation,
		AccountID:     accountId,
		Confirmed:     false}
	err := schema.db.Insert(newTX)
	return *newTX, err
}

func (schema *TransactionSchema) ByTxID(txId string) (Transaction, error) {
	tx := Transaction{}
	err := schema.db.Model(&tx).
		Where("tx_id = ?", txId).
		Select()
	return tx, err
}

func (schema *TransactionSchema) Pending() ([]Transaction, error) {
	var txs []Transaction
	err := schema.db.Model(&txs).
		Where("mined = ?", false).
		Select()
	return txs, err
}

func (schema TransactionSchema) Update(tx Transaction) error {
	err := schema.db.Update(&tx)
	return err
}

func (schema TransactionSchema) UpdateConfirmation(txId string, confirmations int) (error) {
	tx := Transaction{}
	_, err := schema.db.Model(&tx).
		Set("confirmations = ?", confirmations).
		Where("tx_id = ?", txId).
		Returning("*").
		Update()
	if err != nil {
		panic(err)
	}

	var txs []Transaction
	_, err = schema.db.Model(&txs).
		Set("confirmed = ?", true).
		Where("confirmations >= ?", 10).
		Returning("*").
		Update()

	return err
}

func (schema TransactionSchema) Delete(tx Transaction) error {
	err := schema.db.Delete(&tx)
	return err
}
