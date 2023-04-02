package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetEnvVarOrDefault(t *testing.T) {
	tests := []struct {
		envVar       string
		defaultValue string
		expected     string
	}{
		{
			envVar:       "MY_ENV_VAR",
			defaultValue: "default-value",
			expected:     "env-value",
		},
		{
			envVar:       "MY_ENV_VAR",
			defaultValue: "default-value",
			expected:     "default-value",
		},
		{
			envVar:       "NON_EXISTENT_ENV_VAR",
			defaultValue: "default-value",
			expected:     "default-value",
		},
	}

	// Assert it returns the default value as the environment variable is not set
	actual := GetEnvVarOrDefault("MY_ENV_VAR", "default-value")
	assert.Equal(t, "default-value", actual)

	for _, tt := range tests {
		// Set the environment variables for the test cases
		os.Setenv(tt.envVar, tt.expected)

		actual := GetEnvVarOrDefault(tt.envVar, tt.defaultValue)
		assert.Equal(t, tt.expected, actual)
	}

	// Cleanup the environment variables after the test
	os.Unsetenv("MY_ENV_VAR")
}
