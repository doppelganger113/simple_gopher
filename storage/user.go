package storage

import (
	"errors"
	"time"
)

type AuthRole string

const (
	AuthRoleAdmin AuthRole = "Administrators"
	AuthRoleNone  AuthRole = ""
)

func NewAuthRole(value string) (AuthRole, error) {
	converted := AuthRole(value)
	switch converted {
	case AuthRoleAdmin:
		return converted, nil
	case AuthRoleNone:
		return converted, nil
	}

	return "", errors.New("Invalid auth role of " + value)
}

func NewAuthRoleOrDefault(value string, role AuthRole) AuthRole {
	switch AuthRole(value) {
	case AuthRoleAdmin:
		return AuthRoleAdmin
	case AuthRoleNone:
		return AuthRoleNone
	default:
		return role
	}
}

type User struct {
	Id          string     `json:"id"`
	Email       string     `json:"email"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	Role        AuthRole   `json:"role"`
	CogUsername string     `json:"cogUsername"`
	CogSub      string     `json:"cogSub"`
	CogName     string     `json:"cogName"`
	Disabled    bool       `json:"disabled"`
}

type UserList []User

type UserCreationDto struct {
	Email       string
	Role        AuthRole
	CogUsername string
	CogSub      string
	CogName     string
	Disabled    bool
}
