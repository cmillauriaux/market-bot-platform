package model

import (
	"fmt"
	"time"
)

type Statistic struct {
	Date        time.Time
	DisplayDate string
	DateFin     time.Time
	Min         int
	Max         int
	Open        int
	Close       int
	Delta       float64
	Value       int
	Quantity    float64
	Partial     bool
}

func (s *Statistic) Display() string {
	return fmt.Sprintf("%s", s.Date.String()) + "-" + fmt.Sprintf("%s", s.DateFin.String()) + "Min : " + fmt.Sprintf("%v", s.Min) + " Max : " + fmt.Sprintf("%v", s.Max) + " Value : " + fmt.Sprintf("%v", s.Value) + " Quantity : " + fmt.Sprintf("%f", s.Quantity) + " Delta : " + fmt.Sprintf("%f", s.Delta) + "%"
}
