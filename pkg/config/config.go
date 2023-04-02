package config

import "os"

func GetEnvVarOrDefault(envVar string, defaultValue string) string {
	if os.Getenv(envVar) != "" {
		return os.Getenv(envVar)
	}

	return defaultValue
}
