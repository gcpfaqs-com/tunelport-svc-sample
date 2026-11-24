package types

import "errors"

var ErrInvalidName = errors.New("invalid name must be a non-empty string between 5 and 50 characters")

type Name struct {
	value string
}

func NewName(s string) (Name, error) {
	if len(s) < 5 || len(s) > 50 {
		return Name{}, ErrInvalidName
	}

	return Name{
		value: s,
	}, nil
}

func (t *Name) Value() string {
	return t.value
}
