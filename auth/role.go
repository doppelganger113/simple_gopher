package auth

import "errors"

type Role string

const (
	RoleAdmin Role = "Administrators"
	RoleNone  Role = ""
)

func NewAuthRole(value string) (Role, error) {
	converted := Role(value)
	switch converted {
	case RoleAdmin:
		return converted, nil
	case RoleNone:
		return converted, nil
	}

	return "", errors.New("Invalid auth role of " + value)
}
