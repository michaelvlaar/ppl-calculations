package wind

import (
	"fmt"
	"strconv"
)

type Speed float64

func NewSpeedFromString(speed string) (Speed, error) {
	s, err := strconv.ParseFloat(speed, 64)
	if err != nil {
		return Speed(0), err
	}

	return Speed(s), nil
}

func (s Speed) String() string {
	return fmt.Sprintf("%.0f", s)
}
