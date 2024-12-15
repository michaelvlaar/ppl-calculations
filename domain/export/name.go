package export

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidName = errors.New("invalid name")
)

type Name string

func NewName(name string) (Name, error) {
	re := regexp.MustCompile(`^[A-Za-z0-9 ]+$`)
	if !re.MatchString(name) {
		return "", ErrInvalidName
	}
	return Name(name), nil
}

func (name Name) String() string {
	return string(name)
}
