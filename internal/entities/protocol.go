package entities

import "errors"

type Protocol string

const (
	Reality Protocol = "reality"
	NoTLS   Protocol = "notls"
)

func ValidateProtocol(rawProtocol string) error {
	switch Protocol(rawProtocol) {
	case Reality:
		return nil
	case NoTLS:
		return nil
	}
	return errors.New("invalid protocol")
}
