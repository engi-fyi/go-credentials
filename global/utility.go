package global

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// InitTest is used by test functions to initialize logger settings and the assert pkg.
func InitTest(t *testing.T) (*assert.Assertions, zerolog.Logger) {
	os.Setenv(LOG_LEVEL_ENVIRONMENT_KEY, "trace")
	os.Setenv(LOG_OUTPUT_TYPE_ENV_KEY, "pretty")
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return assert.New(t), logger
}
