package serializer

import (
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"github.com/rs/zerolog/log"
	"testing"
)

func TestFromEnv(t *testing.T) {

}

func TestToEnv(t *testing.T) {

}

func TestParseEnvironmentVariable(t *testing.T) {
	assert := global.InitTest(t)
	log.Info().Msg("Testing environment parsing.")
	assert.True(true)

	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	testSerializer := New(testFactory, global.DEFAULT_PROFILE_NAME)

	log.Info().Msg("Testing username parsing.")
	profileName, fieldName, _, didParse := testSerializer.ParseEnvironmentVariable(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL)
	assert.Equal(global.DEFAULT_PROFILE_NAME, profileName)
	assert.Equal(global.TEST_VAR_USERNAME_LABEL, fieldName)
	assert.True(didParse)

	log.Info().Msg("Testing password parsing.")
	profileName, fieldName, _, didParse = testSerializer.ParseEnvironmentVariable(global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL)
	assert.Equal(global.DEFAULT_PROFILE_NAME, profileName)
	assert.Equal(global.TEST_VAR_PASSWORD_LABEL, fieldName)
	assert.True(didParse)

	log.Info().Msg("Testing attribute parsing.")
	profileName, fieldName, sectionName, didParse := testSerializer.ParseEnvironmentVariable(global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL)
	assert.Equal(global.DEFAULT_PROFILE_NAME, profileName)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, fieldName)
	assert.Equal(global.TEST_VAR_FIRST_SECTION_KEY, sectionName)
	assert.True(didParse)


	log.Info().Msg("Testing alternate username parsing.")
	profileName, fieldName, _, didParse = testSerializer.ParseEnvironmentVariable(global.TEST_VAR_ENVIRONMENT_USERNAME_ALTERNATE_LABEL)
	assert.Equal(global.DEFAULT_PROFILE_NAME, profileName)
	assert.Equal(global.TEST_VAR_USERNAME_ALTERNATE_LABEL, fieldName)
	assert.True(didParse)

	log.Info().Msg("Testing alternate password parsing.")
	profileName, fieldName, _, didParse = testSerializer.ParseEnvironmentVariable(global.TEST_VAR_ENVIRONMENT_PASSWORD_ALTERNATE_LABEL)
	assert.Equal(global.DEFAULT_PROFILE_NAME, profileName)
	assert.Equal(global.TEST_VAR_PASSWORD_ALTERNATE_LABEL, fieldName)
	assert.True(didParse)
}