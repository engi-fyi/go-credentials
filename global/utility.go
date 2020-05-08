package global

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestCleanup is a utility function that deletes both the configurationDirectory and credentialFile.
func TestCleanup(configurationDirectory string, credentialFile string) error {
	log.Info().Msg("Cleaning up our global files.")
	if _, credStatErr := os.Stat(credentialFile); !os.IsNotExist(credStatErr) {
		credErr := os.Remove(credentialFile)

		if credErr != nil {
			return credErr
		}
	}

	if _, confStatErr := os.Stat(credentialFile); !os.IsNotExist(confStatErr) {
		confErr := os.Remove(configurationDirectory)

		if confErr != nil {
			return confErr
		}
	}

	return nil
}

// InitTest is used by test functions to initialize logger settings and the assert pkg.
func InitTest(t *testing.T) *assert.Assertions {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return assert.New(t)
}
