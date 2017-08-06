package utils

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/cmillauriaux/market-bot-platform/model"
	uuid "github.com/satori/go.uuid"
)

func ReadCsv(file string, channel chan model.Event, simulationGap time.Duration) error {
	csvfile, err := os.Open(file)

	if err != nil {
		return err
	}

	r := csv.NewReader(csvfile)
	r.Comma = ','
	r.Comment = '#'

	count := 0

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		date, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			return err
		}

		value, err := strconv.ParseFloat(record[1], 64)
		valueInt := int(value * 100)
		if err != nil {
			return err
		}

		quantity, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return err
		}

		tm := time.Unix(date, 0)

		channel <- model.Event{Date: tm, Value: valueInt, Quantity: quantity, OrderID: uuid.NewV4().String()}

		if simulationGap > 0 {
			time.Sleep(simulationGap * time.Millisecond)
		}

		count++
	}

	close(channel)
	return nil
}
