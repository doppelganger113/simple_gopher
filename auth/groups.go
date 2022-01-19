package auth

import (
	"errors"
)

var ErrMissingAuthDto = errors.New("missing auth dto")

func IsInRequiredGroup(cognitoGroups interface{}, requiredGroup Role) bool {
	switch groups := cognitoGroups.(type) {
	case []interface{}:
		for _, group := range groups {
			if group == string(requiredGroup) {
				return true
			}
		}
		return false
	default:
		return false
	}
}
