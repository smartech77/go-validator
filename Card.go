package valitor

import "time"

type VirtualCard struct {
	AquiredAt time.Time
	Number    string
}

type Card struct {
	Test          bool
	Virtual       bool
	CVC           string
	ExpYear       string
	ExpMonth      string
	Number        string
	VirtualNumber string
}

func (c *Card) GetLastFour() string {
	return c.Number[len(c.Number)-5:]
}
