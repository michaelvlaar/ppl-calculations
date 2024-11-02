package callsign

type CallSign string

func New(callSign string) (CallSign, error) {
	return CallSign(callSign), nil
}

func (callSign CallSign) String() string {
	return string(callSign)
}
