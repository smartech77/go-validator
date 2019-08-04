package valitor

import "time"

type VirtualCard struct {
	AquiredAt time.Time
	Number    string
}

type Card struct {
	Test                         bool
	Virtual                      bool
	CVC                          string
	ExpYear                      int
	ExpMonth                     int
	Number                       string
	VirtualNumber                string
	VerificationDataFrom3DSecure string
}

func (c *Card) GetLastFour() string {
	return c.Number[len(c.Number)-5:]
}
