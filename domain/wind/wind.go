package wind

type Wind struct {
	Direction Direction `json:"direction"`
	Speed     Speed     `json:"speed"`
}

func New(direction Direction, speed Speed) (Wind, error) {
	return Wind{Direction: direction, Speed: speed}, nil
}
