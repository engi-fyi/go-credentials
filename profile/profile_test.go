package profile

import (
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
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

func TestRemove(t *testing.T) {
	assert := global.InitTest(t)

	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)

	testProfile := testCreateIniFile(testFactory, assert)
	assert.FileExists(testProfile.ConfigFileLocation)

	removeErr := Remove(testProfile)
	assert.NoError(removeErr)
	assert.Equal(Profile{}, *testProfile)

	os.RemoveAll(testFactory.ParentDirectory)
}

func TestSave(t *testing.T) {
	assert := global.InitTest(t)

	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)

	testProfile := testCreateIniFile(testFactory, assert)
	assert.FileExists(testProfile.ConfigFileLocation)

	os.RemoveAll(testFactory.ParentDirectory)
}

func TestLoad(t *testing.T) {
	assert := global.InitTest(t)

	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)

	log.Info().Msg("Testing the loading of a non-existent profile.")
	testProfile, loadErr := Load(global.TEST_VAR_FIRST_PROFILE_LABEL, testFactory)
	assert.EqualError(loadErr, ERR_PROFILE_DID_NOT_EXIST)
	assert.NoFileExists(testFactory.ConfigDirectory + global.TEST_VAR_FIRST_PROFILE_LABEL)

	testCreateIniFile(testFactory, assert)
	testProfile, loadErr = Load(global.TEST_VAR_FIRST_PROFILE_LABEL, testFactory)
	assert.NoError(loadErr)

	assert.Equal(global.TEST_VAR_NO_SECTION_UNIQUE_KEY_VALUE, testProfile.GetAttribute("", global.TEST_VAR_NO_SECTION_UNIQUE_KEY_LABEL))
	assert.Equal(global.TEST_VAR_DUPLICATE_KEY_VALUE, testProfile.GetAttribute("", global.TEST_VAR_DUPLICATE_KEY_LABEL))
	assert.Equal(global.TEST_VAR_FIRST_SECTION_UNIQUE_KEY_VALUE, testProfile.GetAttribute(global.TEST_VAR_FIRST_SECTION_KEY, global.TEST_VAR_FIRST_SECTION_UNIQUE_KEY_LABEL))
	assert.Equal(global.TEST_VAR_DUPLICATE_KEY_VALUE, testProfile.GetAttribute(global.TEST_VAR_FIRST_SECTION_KEY, global.TEST_VAR_DUPLICATE_KEY_LABEL))
	assert.Equal(global.TEST_VAR_SECOND_SECTION_UNIQUE_KEY_VALUE, testProfile.GetAttribute(global.TEST_VAR_SECOND_SECTION_KEY, global.TEST_VAR_SECOND_SECTION_UNIQUE_KEY_LABEL))
	assert.Equal(global.TEST_VAR_DUPLICATE_KEY_VALUE, testProfile.GetAttribute(global.TEST_VAR_SECOND_SECTION_KEY, global.TEST_VAR_DUPLICATE_KEY_LABEL))
	assert.Empty("", testProfile.GetAttribute(global.TEST_VAR_BAD_SECTION_KEY, ""))
	assert.Empty("", testProfile.GetAttribute("", global.TEST_VAR_BAD_KEY_LABEL))

	os.RemoveAll(testFactory.ParentDirectory)
}

func testCreateIniFile(testFactory *factory.Factory, assert *assert.Assertions) *Profile {
	testProfile, profileErr := New(global.TEST_VAR_FIRST_PROFILE_LABEL, testFactory)
	assert.NoError(profileErr)

	setErr := testProfile.SetAttribute("", global.TEST_VAR_NO_SECTION_UNIQUE_KEY_LABEL, global.TEST_VAR_NO_SECTION_UNIQUE_KEY_VALUE)
	assert.NoError(setErr)
	setErr = testProfile.SetAttribute("", global.TEST_VAR_DUPLICATE_KEY_LABEL, global.TEST_VAR_DUPLICATE_KEY_VALUE)
	assert.NoError(setErr)
	setErr = testProfile.SetAttribute(global.TEST_VAR_FIRST_SECTION_KEY, global.TEST_VAR_FIRST_SECTION_UNIQUE_KEY_LABEL, global.TEST_VAR_FIRST_SECTION_UNIQUE_KEY_VALUE)
	assert.NoError(setErr)
	setErr = testProfile.SetAttribute(global.TEST_VAR_FIRST_SECTION_KEY, global.TEST_VAR_DUPLICATE_KEY_LABEL, global.TEST_VAR_DUPLICATE_KEY_VALUE)
	assert.NoError(setErr)
	setErr = testProfile.SetAttribute(global.TEST_VAR_SECOND_SECTION_KEY, global.TEST_VAR_SECOND_SECTION_UNIQUE_KEY_LABEL, global.TEST_VAR_SECOND_SECTION_UNIQUE_KEY_VALUE)
	assert.NoError(setErr)
	setErr = testProfile.SetAttribute(global.TEST_VAR_SECOND_SECTION_KEY, global.TEST_VAR_DUPLICATE_KEY_LABEL, global.TEST_VAR_DUPLICATE_KEY_VALUE)
	assert.NoError(setErr)
	setErr = testProfile.SetAttribute(global.TEST_VAR_BAD_SECTION_KEY, "", "")
	assert.EqualError(setErr, ERR_MUST_MATCH_REGEX)
	setErr = testProfile.SetAttribute("", global.TEST_VAR_BAD_KEY_LABEL, "")
	assert.EqualError(setErr, ERR_MUST_MATCH_REGEX)

	saveErr := testProfile.Save()
	assert.NoError(saveErr)
	assert.FileExists(testProfile.ConfigFileLocation)

	return testProfile
}