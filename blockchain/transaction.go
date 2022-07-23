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

func (c *Chain) AddTransaction(sender string, recp string, amount float64) string {
	t := Transaction{
		Sender:    sender,
		Recipient: recp,
		Amount:    amount,
	}
	jsonString, _ := json.Marshal(t)

	return c.AddBlock(string(jsonString))
}

func (t *Transaction) String() string {
	return fmt.Sprintf("%s -> %s: %.2f", t.Sender, t.Recipient, t.Amount)
}

func (t *Transaction) Validate() bool {
	if t.Sender == "" || t.Recipient == "" || t.Amount <= 0 {
		return false
	}
	return true
}
