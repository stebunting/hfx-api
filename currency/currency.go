package currency

import (
	"errors"
	"regexp"
)

type Code [3]byte

func NewCode(code string) (Code, error) {
	var c Code

	valid, err := regexp.MatchString("^[A-Z]{3}$", code)
	if err != nil {
		panic(err)
	}

	if !valid {
		return c, errors.New("invalid currency code")
	}
	copy(c[:], code)

	return c, nil
}

func (c Code) String() string {
	return string(c[:])
}
