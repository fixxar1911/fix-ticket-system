package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "test_value")
	defer os.Unsetenv("TEST_ENV_VAR")

	// Should return the set value
	val := getEnv("TEST_ENV_VAR", "default")
	assert.Equal(t, "test_value", val)

	// Should return the default value if not set
	val = getEnv("NON_EXISTENT_ENV_VAR", "default")
	assert.Equal(t, "default", val)
}

func TestInitDB_Skip(t *testing.T) {
	t.Skip("Skipping InitDB test as it requires a real Postgres instance or mocking.")
}
