package profile

import (
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	assert := global.InitTest(t)

	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)

	testProfile, newErr := New(global.TEST_VAR_FIRST_PROFILE_LABEL, testFactory)
	assert.NoError(newErr)
	assert.Equal(global.TEST_VAR_FIRST_PROFILE_LABEL, testProfile.Name)
	assert.Equal(testFactory.ConfigDirectory+global.TEST_VAR_FIRST_PROFILE_LABEL, testProfile.ConfigFileLocation)
	assert.Equal(testProfile.Factory, testFactory)
	assert.True(testProfile.Initialized)

	testProfile, newErr = New(global.TEST_VAR_BAD_PROFILE_LABEL, testFactory)
	assert.EqualError(newErr, ERR_PROFILE_NAME_MUST_MATCH_REGEX)
	assert.Nil(testProfile)

	os.RemoveAll(testFactory.ParentDirectory)
}

func TestAttributes(t *testing.T) {

}