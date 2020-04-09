package models

import "time"

type Session struct {
	SID string
	ID uint
	Expires time.Time
}
