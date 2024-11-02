package pressure

import (
	"fmt"
	"strconv"
)

type Altitude float64

func NewFromString(s string) (Altitude, error) {
	a, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return Altitude(0.0), err
	}

	return Altitude(a), nil
}

func (a Altitude) Float64() float64 {
	return float64(a)
}

func (a Altitude) String() string {
	return fmt.Sprintf("%.0f", a)
}
