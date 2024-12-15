package export

import (
	"github.com/google/uuid"
)

type ID uuid.UUID

func NewID() (ID, error) {
	return ID(uuid.New()), nil
}

func NewIDFromString(s string) (ID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return ID(uuid.UUID{}), err
	}
	return ID(u), nil
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}
