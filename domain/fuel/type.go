package fuel

type Type string

const (
	TypeMogas Type = "mogas"
	TypeAvGas Type = "avgas"
)

func NewTypeFromString(s string) (Type, error) {
	if s == "mogas" {
		return TypeMogas, nil
	} else {
		return TypeAvGas, nil
	}
}

func (t Type) String() string {
	return string(t)
}
