package seat

import (
	"errors"
)

type Position int

const (
	PositionFront Position = iota
	PositionMiddle
	PositionBack
)

var ErrInvalidPosition = errors.New("invalid position")

func NewFromString(position string) (Position, error) {
	if position == "f" {
		return PositionFront, nil
	} else if position == "m" {
		return PositionMiddle, nil
	} else if position == "b" {
		return PositionBack, nil
	}

	return Position(-1), ErrInvalidPosition
}

func (p Position) String() string {
	if p == PositionFront {
		return "f"
	} else if p == PositionMiddle {
		return "m"
	} else if p == PositionBack {
		return "b"
	}

	panic("invalid position")
}
