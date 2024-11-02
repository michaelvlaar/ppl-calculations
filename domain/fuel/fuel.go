package fuel

import "ppl-calculations/domain/volume"

type Fuel struct {
	Volume volume.Volume `json:"volume"`
	Type   Type          `json:"type"`
}

func New(volume volume.Volume, t Type) (Fuel, error) {
	return Fuel{Volume: volume, Type: t}, nil
}

func MustNew(volume volume.Volume, t Type) Fuel {
	f, err := New(volume, t)
	if err != nil {
		panic(err)
	}

	return f
}

func Add(fuels ...Fuel) Fuel {
	v := volume.MustNew(0, volume.TypeLiter)
	for _, f := range fuels {
		if fuels[0].Type != f.Type {
			panic("fuel types must match")
		}

		v = v.Add(f.Volume)
	}

	return MustNew(v, fuels[0].Type)
}
