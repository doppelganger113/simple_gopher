package http_util

import (
	"strconv"
	"strings"
)

func ToUint(value string) uint {
	if value == "" {
		return 0
	}
	res, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0
	}

	return uint(res)
}

func GetTokenFromHeader(authHeader string) string {
	parts := strings.SplitAfter(authHeader, "Bearer ")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}
