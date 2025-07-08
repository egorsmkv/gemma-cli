package env

import (
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// FatalOnMissingEnv allows the app to panic if the `key` variable is not found
var FatalOnMissingEnv bool

// LoadFromFile - loads .env file by name
func LoadFromFile(path string) error {
	return godotenv.Overload(path)
}

// GetAsString gets a string from the environment
func GetAsString(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists && FatalOnMissingEnv {
		slog.Error("Missing environment variable", "key", key)
		os.Exit(1)
	}
	return value
}

// GetAsStringElse gets a string from the environment and if not found, uses the alt
func GetAsStringElse(key string, alt string) string {
	if FatalOnMissingEnv {
		panic("env.FatalOnMissingEnv is incompatible with using *Else() functions")
	}
	value, exists := os.LookupEnv(key)
	if !exists {
		return alt
	}
	return value
}

// GetAsInt gets an integer from the environment
func GetAsInt(key string) int {
	valueStr := GetAsString(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil && FatalOnMissingEnv {
		slog.Error("Environment variable is not an integer", "key", key)
		os.Exit(1)
	}
	return value
}

// GetAsIntElse gets an integer from the environment and if not found, uses the alt
func GetAsIntElse(key string, alt int) int {
	if FatalOnMissingEnv {
		panic("env.FatalOnMissingEnv is incompatible with using *Else() functions")
	}
	valueStr := GetAsString(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return alt
	}
	return value
}

// GetAsBool gets a boolean from the environment
func GetAsBool(key string) bool {
	valueStr := GetAsString(key)
	value, err := strconv.ParseBool(valueStr)
	if err != nil && FatalOnMissingEnv {
		slog.Error("Environment variable is not a boolean", "key", key)
		os.Exit(1)
	}
	return value
}

// GetAsBoolElse gets a boolean from the environment and if not found, uses the alt
func GetAsBoolElse(key string, alt bool) bool {
	if FatalOnMissingEnv {
		panic("env.FatalOnMissingEnv is incompatible with using *Else() functions")
	}
	valueStr := GetAsString(key)
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return alt
	}
	return value
}

// GetAsSlice gets a string from the environment and splits it into a slice using the specified separator
func GetAsSlice(name string, sep string) []string {
	valStr := GetAsString(name)
	return strings.Split(valStr, sep)
}

// GetAsSliceElse gets a string from the environment, splits it into a slice using the specified separator, and if not found, uses the alt
func GetAsSliceElse(name string, sep string, alt []string) []string {
	if FatalOnMissingEnv {
		panic("env.FatalOnMissingEnv is incompatible with using *Else() functions")
	}
	valStr := GetAsString(name)
	if len(valStr) == 0 {
		return alt
	}
	return strings.Split(valStr, sep)
}
