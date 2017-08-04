package utils

import "time"

type Counter struct {
	begin time.Time
}

func (c *Counter) StartCount() {
	c.begin = time.Now()
}

func (c *Counter) StopCount() time.Duration {
	return time.Since(c.begin)
}
