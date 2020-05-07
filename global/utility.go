package global

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCleanup(configurationDirectory string, credentialFile string) error {
	log.Info().Msg("Cleaning up our global files.")
	credErr := os.Remove(credentialFile)
	confErr := os.Remove(configurationDirectory)

	if confErr != nil {
		return credErr
	}

	if credErr != nil {
		return credErr
	}

	return nil
}

func InitTest(t *testing.T) *assert.Assertions {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return assert.New(t)
}
