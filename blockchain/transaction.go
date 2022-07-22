package blockchain

import "encoding/json"

type Transaction struct {
	Sender    string
	Recipient string
	Amount    float64
}

func (c *Chain) AddTransaction(sender string, recp string, amount float64) {
	t := Transaction{
		Sender:    sender,
		Recipient: recp,
		Amount:    amount,
	}
	jsonString, _ := json.Marshal(t)

	c.AddBlock(string(jsonString))
}
