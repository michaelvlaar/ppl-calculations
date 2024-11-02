package wind

import "errors"

type Direction int

const (
	DirectionTailwind Direction = iota
	DirectionHeadwind
)

var (
	ErrInvalidDirection = errors.New("Invalid direction")
)

func NewDirectionFromString(s string) (Direction, error) {
	if s == "headwind" {
		return DirectionHeadwind, nil
	} else if s == "tailwind" {
		return DirectionTailwind, nil
	}

	return -1, ErrInvalidDirection
}

func (d Direction) String() string {
	if d == DirectionHeadwind {
		return "headwind"
	} else if d == DirectionTailwind {
		return "tailwind"
	}

	panic("Invalid direction")
}
