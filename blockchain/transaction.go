package blockchain

import (
	"encoding/json"
	"fmt"
)

type Transaction struct {
	Sender    string
	Recipient string
	Amount    float64
}

func NewTransactionBlock(sender string, recp string, amount float64) (*Block, error) {
	t := Transaction{
		Sender:    sender,
		Recipient: recp,
		Amount:    amount,
	}

	jsonString, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return NewBlock(string(jsonString)), nil
}

// String representation of a transaction
func (t *Transaction) String() string {
	return fmt.Sprintf("%s -> %s: %.2f", t.Sender, t.Recipient, t.Amount)
}

// Check transaction validity
func (t *Transaction) Validate() bool {
	if t.Sender == "" || t.Recipient == "" || t.Amount <= 0 {
		return false
	}

	return true
}
