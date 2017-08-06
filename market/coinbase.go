package market

import (
	"log"
	"os"
	"strconv"
	"time"

	"git.icysoft.fr/cedric/kraken-bot/model"
	ws "github.com/gorilla/websocket"
	exchange "github.com/preichenberger/go-coinbase-exchange"
)

type Coinbase struct {
	client *exchange.Client
}

func InitMarket() Market {
	secret := os.Getenv("COINBASE_SECRET")
	key := os.Getenv("COINBASE_KEY")
	passphrase := ""

	coinbase := Coinbase{client: exchange.NewClient(secret, key, passphrase)}

	return &coinbase
}

func (c *Coinbase) GetStatistic(start time.Time, end time.Time) (*model.Statistic, error) {
	params := exchange.GetHistoricRatesParams{Start: start, End: end, Granularity: 60 * 60 * 24}
	rates, err := c.client.GetHistoricRates("BTC-EUR", params)

	if err != nil {
		return nil, err
	}

	totalQuantity := 0.0
	totalValue := 0
	nbValues := 0
	min := 0
	max := 0

	for _, rate := range rates {
		totalQuantity += rate.Volume
		nbValues++
		if min == 0 || int(rate.Low*100) < min {
			min = int(rate.Low * 100)
		}
		if max == 0 || int(rate.High*100) > max {
			max = int(rate.High * 100)
		}
	}

	delta := 0.0
	if min > 0 {
		delta = (float64(max-min) / float64(min) * 100)
	}

	value := 0
	if nbValues > 0 {
		value = totalValue / nbValues
	}

	return &model.Statistic{Min: min, Max: max, Quantity: totalQuantity, Value: value, Delta: delta, Date: start, DateFin: end}, nil
}

func (c *Coinbase) GetTransactions(start time.Time, end time.Time) ([]*model.Event, error) {
	params := exchange.ListTradesParams{}
	cursor := c.client.ListTrades("BTC-EUR", params)
	var trades []exchange.Trade
	events := make([]*model.Event, 0)

	for cursor.HasMore {
		if err := cursor.NextPage(&trades); err != nil {
			return nil, err
		}

		for _, o := range trades {
			if o.Time.Time().After(start) && o.Time.Time().Before(end) {
				events = append(events, &model.Event{OrderID: strconv.FormatInt(int64(o.TradeId), 10), Date: o.Time.Time(), Quantity: o.Size, Value: int(o.Price * 100)})
			} else {
				return events, nil
			}
		}
	}
	return events, nil
}

func (c *Coinbase) GetSleepTimeBetweenRequests() time.Duration {
	return time.Second
}

func (c *Coinbase) SubscribeToFlux(broacastFn BroadcastEvent) {
	go c.connectToWebservice(broacastFn)
}

func (c *Coinbase) connectToWebservice(broacastFn BroadcastEvent) {
	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial("wss://ws-feed.gdax.com", nil)
	if err != nil {
		println(err.Error())
	}

	subscribe := map[string]string{
		"type":       "subscribe",
		"product_id": "BTC-EUR",
	}
	if err := wsConn.WriteJSON(subscribe); err != nil {
		println(err.Error())
	}

	message := exchange.Message{}
	for true {
		if err := wsConn.ReadJSON(&message); err != nil {
			println(err.Error())
		}

		if message.Type == "match" {
			event := model.Event{OrderID: message.OrderId, Quantity: message.Size, Value: int(message.Price * 100), Date: message.Time.Time()}
			log.Println(event)
			broacastFn(&event)
		}
	}
}
