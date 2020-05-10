package global

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// InitTest is used by test functions to initialize logger settings and the assert pkg.
func InitTest(t *testing.T) *assert.Assertions {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return assert.New(t)
}
