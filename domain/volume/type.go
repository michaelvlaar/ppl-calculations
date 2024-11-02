package volume

type Type string

const (
	TypeLiter  Type = "liter"
	TypeGallon Type = "gallon"
)

func NewTypeFromString(s string) (Type, error) {
	if s == "liter" {
		return TypeLiter, nil
	} else {
		return TypeGallon, nil
	}
}

func (t Type) String() string {
	return string(t)
}
