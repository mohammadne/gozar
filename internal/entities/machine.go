package entities

import "errors"

type Machine string

const (
	Client Machine = "client"
	// Middleware Machine = "middleware"
	Server Machine = "server"
)

func ValidateMachine(rawMachine string) error {
	switch Machine(rawMachine) {
	case Client:
		return nil
	case Server:
		return nil
	}
	return errors.New("invalid machine type")
}
