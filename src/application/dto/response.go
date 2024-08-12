package dto

import "time"

type Response struct {
	Timestamp time.Time
	Object    any
	Messages  []string
}
