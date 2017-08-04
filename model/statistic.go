package model

import (
	"fmt"
	"time"
)

type Statistic struct {
	Date     time.Time
	Min      int
	Max      int
	Delta    float64
	Value    int
	Quantity float64
}

func (s *Statistic) Display() string {
	return "DATE : " + fmt.Sprintf("%s", s.Date.String()) + "Min : " + fmt.Sprintf("%v", s.Min) + " Max : " + fmt.Sprintf("%v", s.Max) + " Value : " + fmt.Sprintf("%v", s.Value) + " Quantity : " + fmt.Sprintf("%f", s.Quantity) + " Delta : " + fmt.Sprintf("%f", s.Delta) + "%"
}
