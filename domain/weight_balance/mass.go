package weight_balance

import (
	"fmt"
	"strconv"
)

type Mass float64

func NewMassFromString(weight string) (Mass, error) {
	w, err := strconv.ParseFloat(weight, 64)
	if err != nil {
		return Mass(0), err
	}

	return Mass(w), nil
}

func NewMass(mass float64) Mass {
	return Mass(mass)
}

func (w *Mass) String() string {
	return fmt.Sprintf("%.0f", *w)
}

func (w *Mass) Kilo() float64 {
	return float64(*w)
}

type MassMoment struct {
	Name string
	// Arm in meter
	Arm float64
	// Mass in KG
	Mass Mass
}

func NewMassMoment(name string, arm float64, mass Mass) *MassMoment {
	return &MassMoment{Name: name, Arm: arm, Mass: mass}
}

func (mm *MassMoment) KGM() float64 {
	return mm.Arm * mm.Mass.Kilo()
}
