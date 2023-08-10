package config

import (
	"errors"
	"time"
)

var ErrZeroInterval = errors.New("zero interval not allowed")

type Poller struct {
	Interval time.Duration
}

func (p *Poller) load(envPrefix string) error {
	v := setupViper(envPrefix)

	p.Interval = v.GetDuration("interval")
	if p.Interval == 0 {
		return ErrZeroInterval
	}

	return nil
}
