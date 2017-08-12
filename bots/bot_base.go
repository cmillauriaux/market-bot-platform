package bots

import (
	"log"

	"github.com/cmillauriaux/market-bot-platform/market"
	"github.com/cmillauriaux/market-bot-platform/model"
	uuid "github.com/satori/go.uuid"
)

type BotBase struct {
	ID            string
	client        market.Market
	Transactions  map[string]*model.Event `json:"-"`
	Orders        map[string]*model.Order `json:"-"`
	ReverseOrders map[string]*model.Order `json:"-"`
	//History      map[string]*core.Transaction `json:"-"`
	NbTransactions int
	Wallet         int
}

func (b *BotBase) Init(client market.Market, wallet int) {
	b.ID = uuid.NewV4().String()
	b.client = client
	b.Wallet = wallet
	b.Transactions = make(map[string]*model.Event)
	b.Orders = make(map[string]*model.Order)
	b.ReverseOrders = make(map[string]*model.Order)
	//b.History = make(map[string]*core.Transaction)
}

func (b *BotBase) BuySuccess(transaction *model.Event, order *model.Order) {
	// Make transaction
	b.Transactions[transaction.OrderID] = transaction
	b.Transactions[transaction.OrderID].Quantity = order.Quantity

	// Remove order
	delete(b.Orders, order.OrderID)
	delete(b.ReverseOrders, order.TransactionID)

	// Substract from Wallet
	b.Wallet -= int(b.Transactions[transaction.OrderID].Quantity * float64(b.Transactions[transaction.OrderID].Value))
	//log.Println("BUY SUCCESS [", b.ID, ":", order.OrderID, "] : Order [", order.Value, "] Tx [", transaction.Value, "] Qty [", b.Transactions[transaction.OrderID].Quantity, "] VAL [", float64(b.Transactions[transaction.OrderID].Quantity*float64(b.Transactions[transaction.OrderID].Value)), "]")
	//log.Println("WALLET STATE [", b.ID, "]: ", b.Wallet)

	// Register History
	b.NbTransactions++
	//b.History[transaction.OrderID] = transaction
}

func (b *BotBase) SellSuccess(transaction *model.Event, order *model.Order) {
	// Remove order
	delete(b.Orders, order.OrderID)
	delete(b.ReverseOrders, order.TransactionID)

	// Remove transaction
	delete(b.Transactions, order.TransactionID)

	// Increment wallet
	b.Wallet += int(order.Quantity * float64(transaction.Value))
	//log.Println("SELL SUCCESS [", b.ID, ":", order.OrderID, "] : Order [", order.Value, "] Tx [", transaction.Value, "] Qty [", order.Quantity, "] VAL [", float64(order.Quantity*float64(transaction.Value)), "]")
	//log.Println("WALLET STATE [", b.ID, "]: ", b.Wallet)

	// Register History
	transaction.Quantity = order.Quantity
	//b.History[transaction.OrderID] = transaction
	b.NbTransactions++
}

func (b *BotBase) GetWalletValue() int {
	return b.Wallet - b.GetTotalOrdersValue()
}

func (b *BotBase) GetTotalOrdersValue() int {
	total := 0
	for _, transaction := range b.Orders {
		if transaction.Buy {
			total += int(float64(transaction.Value) * transaction.Quantity)
		}
	}
	return total
}

func (b *BotBase) GetTotalTransactionValue() int {
	total := 0
	for _, transaction := range b.Transactions {
		total += int(float64(transaction.Value) * transaction.Quantity)
	}
	return total
}

func (b *BotBase) Display() {
	log.Println("[", b.ID, "] Wallet : ", b.Wallet, " (", b.Wallet+b.GetTotalTransactionValue(), ") orders : ", len(b.Orders), " tx : ", len(b.Transactions), " nbTx : ", b.NbTransactions)
}
