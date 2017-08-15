package market

import (
	"sync"
	"time"

	"github.com/cmillauriaux/market-bot-platform/model"
	uuid "github.com/satori/go.uuid"
)

type TrainingClient struct {
	orders map[string]*model.Order
	mutex  *sync.Mutex
}

func (c *TrainingClient) GetStatistic(start time.Time, end time.Time) (*model.Statistic, error) {
	return &model.Statistic{}, nil
}

func (c *TrainingClient) GetTransactions(start time.Time, end time.Time) ([]*model.Event, error) {
	events := make([]*model.Event, 0)
	return events, nil
}

func (c *TrainingClient) GetSleepTimeBetweenRequests() time.Duration {
	return time.Nanosecond
}

func (c *TrainingClient) SubscribeToFlux(broacastFn BroadcastEvent) {
}

func (c *TrainingClient) SimulateMarketTransaction(event *model.Event) {
	if c.orders == nil {
		c.orders = make(map[string]*model.Order)
	}
	for key, order := range c.orders {
		if event.Value <= order.OrderValue && order.Buy && order.BuySuccess != nil {
			order.BuySuccess(event, order)
			delete(c.orders, key)
		}
		if event.Value >= order.OrderValue && order.Sell && order.SellSuccess != nil {
			order.SellSuccess(event, order)
			delete(c.orders, key)
		}
	}
}

func (c *TrainingClient) CancelOrder(orderId string) {
	if c.mutex == nil {
		c.mutex = &sync.Mutex{}
	}
	c.mutex.Lock()
	delete(c.orders, orderId)
	c.mutex.Unlock()
}

func (c *TrainingClient) MakeBuyOrder(size float64, value int, originalValue int, callback func(*model.Event, *model.Order)) *model.Order {
	if c.orders == nil {
		c.orders = make(map[string]*model.Order)
	}
	if c.mutex == nil {
		c.mutex = &sync.Mutex{}
	}
	c.mutex.Lock()
	transaction := model.Order{}
	transaction.OrderID = uuid.NewV4().String()
	transaction.Buy = true
	transaction.Quantity = size
	transaction.OrderValue = value
	transaction.OriginalValue = originalValue
	transaction.BuySuccess = callback
	c.orders[transaction.OrderID] = &transaction
	c.mutex.Unlock()
	return &transaction
}

func (c *TrainingClient) MakeSellOrder(size float64, value int, originalValue int, callback func(*model.Event, *model.Order)) *model.Order {
	if c.orders == nil {
		c.orders = make(map[string]*model.Order)
	}
	if c.mutex == nil {
		c.mutex = &sync.Mutex{}
	}
	c.mutex.Lock()
	transaction := model.Order{}
	transaction.OrderID = uuid.NewV4().String()
	transaction.Sell = true
	transaction.Quantity = size
	transaction.OrderValue = value
	transaction.OriginalValue = originalValue
	transaction.SellSuccess = callback
	c.orders[transaction.OrderID] = &transaction
	c.mutex.Unlock()
	return &transaction
}
