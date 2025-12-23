package utils

import (
	"os"
)

func GetEnv(key string, defaultVal ...string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return ""
}
