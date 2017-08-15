package bots

import (
	"log"
	"sort"

	"github.com/cmillauriaux/market-bot-platform/history"

	"github.com/cmillauriaux/market-bot-platform/market"
	"github.com/cmillauriaux/market-bot-platform/model"
	uuid "github.com/satori/go.uuid"
)

type BotBase struct {
	ID             string
	client         market.Market
	Statistics     *history.History
	Transactions   map[string]*model.Order `json:"-"`
	Orders         map[string]*model.Order `json:"-"`
	History        map[string]*model.Order `json:"-"`
	NbTransactions int
	Wallet         int
}

func (b *BotBase) GetID() string {
	return b.ID
}

func (b *BotBase) Init(client market.Market, statistics *history.History, wallet int) {
	b.ID = uuid.NewV4().String()
	b.client = client
	b.Wallet = wallet
	b.Transactions = make(map[string]*model.Order)
	b.Orders = make(map[string]*model.Order)
	b.History = make(map[string]*model.Order)
	b.Statistics = statistics
}

func (b *BotBase) BuySuccess(event *model.Event, order *model.Order) {
	// Remove order
	if b.Orders[order.OrderID] == nil {
		log.Fatal("CAN'T SELL IT [ORDER]")
	}
	delete(b.Orders, order.OrderID)

	if b.Transactions[order.OrderID] == nil {
		// Register History
		b.NbTransactions++
		order.TransactionID = event.OrderID
		order.Buy = true
		order.Sell = false
		order.DisplayDate = event.Date.Format("2006-01-02 15:04:05")
		order.TransactionValue = event.Value
		b.History[order.OrderID] = order

		// Register transaction
		b.Transactions[order.OrderID] = order

		// Substract from Wallet
		b.Wallet -= order.GetTransactionValue()
		//log.Println("BUY SUCCESS ", event.Date, " Order [", order.OrderValue, "] Tx [", order.TransactionValue, "] Qty [", order.Quantity, "] VAL [", order.GetTransactionValue(), "] Orders :", " Wallet : ", b.Wallet)
		//log.Println("WALLET STATE [", b.ID, "]: ", b.Wallet)
	} else {
		log.Fatal("CAN'T BUY IT")
	}
}

func (b *BotBase) SellSuccess(event *model.Event, order *model.Order) {
	// Remove order
	if b.Orders[order.OrderID] == nil {
		log.Fatal("CAN'T SELL IT [ORDER]")
	}
	delete(b.Orders, order.OrderID)

	isTxFound := 0
	for transactionID, transaction := range b.Transactions {
		if transaction.TransactionID == order.TransactionID {
			delete(b.Transactions, transactionID)
			isTxFound++
		}
	}
	if isTxFound != 1 {
		//log.Fatal("TX NOT FOUND : ", isTxFound)
		log.Println("TX NOT FOUND : ", isTxFound)
	}

	// Register History
	order.Buy = false
	order.Sell = true
	order.DisplayDate = event.Date.Format("2006-01-02 15:04:05")
	order.TransactionValue = event.Value
	b.History[order.OrderID] = order
	b.NbTransactions++

	// Increment wallet
	b.Wallet += order.GetTransactionValue()
	//log.Println("SELL FOR : ", order.GetTransactionValue(), order.GetOriginalValue(), order.GetPlusValue())
	//log.Println("SELL SUCCESS ", event.Date, " Order [", order.OrderValue, "] Tx [", event.Value, "] Qty [", order.Quantity, "] VAL [", float64(order.Quantity*float64(event.Value)), "]")
	//log.Println("WALLET STATE [", b.ID, "]: ", b.Wallet)
}

func (b *BotBase) GetWalletValue() int {
	return b.Wallet - b.GetTotalOrdersValue()
}

func (b *BotBase) GetTotalOrdersValue() int {
	total := 0
	for _, transaction := range b.Orders {
		if transaction.Buy {
			total += transaction.GetOriginalValue()
		}
	}
	return total
}

func (b *BotBase) GetTotalTransactionValue() int {
	total := 0
	for _, transaction := range b.Transactions {
		total += transaction.GetTransactionValue()
	}
	return total
}

func (b *BotBase) GetTotalValue() int {
	return b.Wallet + b.GetTotalTransactionValue() + b.GetTotalOrdersValue()
}

func (b *BotBase) GetHistory() []*model.Order {
	history := make([]*model.Order, 0)
	for _, order := range b.History {
		history = append(history, order)
	}
	sort.SliceStable(history, func(i, j int) bool { return history[i].Date.Before(history[j].Date) })
	return history
}
