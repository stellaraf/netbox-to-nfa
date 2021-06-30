package util

import (
	"fmt"
	"os"
)

// CheckEnv ensures an environment variable exists and is not empty, if it is required.
func CheckEnv(key string, required bool) {
	value := os.Getenv(key)
	if required && value == "" {
		fmt.Printf("Unable to find required environment variable '%s'\n", key)
		os.Exit(1)
	}
}

// GetEnv gets an environment variable value. This will likely expand to include fanciness.
func GetEnv(key string) (value string) {
	return os.Getenv(key)
}
