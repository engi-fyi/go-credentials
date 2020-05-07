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
	os.Remove(credentialFile)
	os.Remove(configurationDirectory)

	return nil
}

func InitTest(t *testing.T) *assert.Assertions {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return assert.New(t)
}
