package utils

import (
	"encoding/csv"
	"os"
)

func WriteCsv(fileName string, data [][]string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			return err
		}
	}

	defer writer.Flush()
	return nil
}
