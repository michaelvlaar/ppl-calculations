package weight

import (
	"fmt"
	"strconv"
)

type Weight float64

func NewFromString(weight string) (Weight, error) {
	w, err := strconv.ParseFloat(weight, 64)
	if err != nil {
		return Weight(0), err
	}

	return Weight(w), nil
}

func (w *Weight) String() string {
	return fmt.Sprintf("%.0f", *w)
}
