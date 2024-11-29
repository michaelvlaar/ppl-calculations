package volume

import "fmt"

const (
	LitersInGallon = 3.78541

	DescriptionLiter  = "L"
	DescriptionGallon = "gal"
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

func (v Volume) Subtract(other Volume) Volume {
	return Subtract(v, other)
}

func Subtract(base Volume, volumes ...Volume) Volume {
	var amount float64
	switch base.Type {
	case TypeLiter:
		amount = base.Amount
	case TypeGallon:
		amount = base.Amount * LitersInGallon
	}

	for _, v := range volumes {
		switch v.Type {
		case TypeLiter:
			amount -= v.Amount
		case TypeGallon:
			amount -= v.Amount * LitersInGallon
		}
	}

	return MustNew(amount, TypeLiter)
}
func (v Volume) String(t Type) string {
	var amount float64
	switch v.Type {
	case TypeLiter:
		amount = v.Amount
	case TypeGallon:
		amount = v.Amount * LitersInGallon
	}

	switch t {
	case TypeLiter:
		return fmt.Sprintf("%.2f%s", amount, DescriptionLiter)
	case TypeGallon:
		return fmt.Sprintf("%.2f%s", amount/LitersInGallon, DescriptionGallon)
	default:
		panic("invalid volume type")
	}
}
