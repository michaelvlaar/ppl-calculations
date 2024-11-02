package temperature

import (
	"fmt"
	"strconv"
)

type Temperature float64

func NewFromString(temperature string) (Temperature, error) {
	t, err := strconv.ParseFloat(temperature, 64)
	if err != nil {
		return Temperature(0), err
	}

	return Temperature(t), nil
}

func (t Temperature) String() string {
	return fmt.Sprintf("%.0f", float64(t))
}
