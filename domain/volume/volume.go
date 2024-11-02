package volume

import "fmt"

const (
	LitersInGallon = 3.78541
)

type Volume struct {
	Amount float64
	Type   Type
}

func New(amount float64, t Type) (Volume, error) {
	return Volume{amount, t}, nil
}

func MustNew(amount float64, t Type) Volume {
	v, err := New(amount, t)
	if err != nil {
		panic(err)
	}

	return v
}

func (v Volume) Add(other Volume) Volume {
	return Add(v, other)
}

func Add(volumes ...Volume) Volume {
	amount := 0.0
	for _, v := range volumes {
		switch v.Type {
		case TypeLiter:
			amount += v.Amount
		case TypeGallon:
			amount += v.Amount * LitersInGallon
		}
	}

	return MustNew(amount, TypeLiter)
}

func (v Volume) String() string {
	switch v.Type {
	case TypeLiter:
		return fmt.Sprintf("%.2fL", v.Amount)
	case TypeGallon:
		return fmt.Sprintf("%.2fgal", v.Amount)
	default:
		panic("invalid volume type")
	}
}
