package models

import (
	"fmt"
)

type Segment struct {
	NameSegment string `json:"name" validate:"required"`
}

type SubscribeRequest struct {
	UserId       int    `json:"user-id" validate:"required"`
	SegmentName  string `json:"segment-name" validate:"required"`
	TimeoutHours int    `json:"timeout-hours"`
}

type Respond struct {
	Status  int
	Message string
	Err     error
}

func (r *Respond) Error() string {
	return fmt.Sprintf("Status: %d\nMessage: %s\nError: %v",
		r.Status, r.Message, r.Err)
}
