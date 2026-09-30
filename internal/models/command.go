package models

import "time"

type Command struct {
	Arguments []string
	CWD       string
	Timeout   time.Duration
}
