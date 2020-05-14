package profile

import (
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"os"
	"testing"
)

func TestProfileNew(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing the creation of a new profile.")
	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(factoryErr)

	testProfile, newErr := New(global.TEST_VAR_FIRST_PROFILE_LABEL, testFactory)
	assert.NoError(newErr)
	assert.Equal(global.TEST_VAR_FIRST_PROFILE_LABEL, testProfile.Name)
	assert.Equal(testFactory.ConfigDirectory+global.TEST_VAR_FIRST_PROFILE_LABEL, testProfile.ConfigFileLocation)
	assert.Equal(testProfile.Factory, testFactory)
	assert.True(testProfile.Initialized)

	os.RemoveAll(testFactory.ParentDirectory)
}

func TestProfileBadName(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing the creation of a new profile with a bad name.")
	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(factoryErr)
	testProfile, newErr := New(global.TEST_VAR_BAD_PROFILE_LABEL, testFactory)
	assert.EqualError(newErr, ERR_PROFILE_NAME_MUST_MATCH_REGEX)
	assert.Nil(testProfile)
}

func TestProfileBlankName(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing the creation of a new profile with a bad name.")
	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(factoryErr)
	testProfile, newErr := New("", testFactory)
	assert.EqualError(newErr, ERR_PROFILE_NAME_MUST_MATCH_REGEX)
	assert.Nil(testProfile)
}

func TestProfileAttribute(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing the single attributes on a profile.")
	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(factoryErr)

	testProfile, newErr := New(global.TEST_VAR_FIRST_PROFILE_LABEL, testFactory)
	assert.NoError(newErr)
	assert.Equal(global.TEST_VAR_FIRST_PROFILE_LABEL, testProfile.Name)
	assert.Equal(testFactory.ConfigDirectory+global.TEST_VAR_FIRST_PROFILE_LABEL, testProfile.ConfigFileLocation)
	assert.Equal(testProfile.Factory, testFactory)
	assert.True(testProfile.Initialized)

	setErr := testProfile.SetAttribute(global.NO_SECTION_KEY, global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.NoError(setErr)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, testProfile.GetAttribute(global.NO_SECTION_KEY, global.TEST_VAR_ATTRIBUTE_NAME_LABEL))

	os.RemoveAll(testFactory.ParentDirectory)
}

func TestProfileAllAttributes(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing getting all attributes from a profile.")
	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(factoryErr)

	testProfile, newErr := New(global.TEST_VAR_FIRST_PROFILE_LABEL, testFactory)
	assert.NoError(newErr)
	assert.Equal(global.TEST_VAR_FIRST_PROFILE_LABEL, testProfile.Name)
	assert.Equal(testFactory.ConfigDirectory+global.TEST_VAR_FIRST_PROFILE_LABEL, testProfile.ConfigFileLocation)
	assert.Equal(testProfile.Factory, testFactory)
	assert.True(testProfile.Initialized)

	setErr := testProfile.SetAttribute(global.SECTION_NAME_BLANK, global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.NoError(setErr)
	setErr = testProfile.SetAttribute(global.TEST_VAR_FIRST_SECTION_KEY, global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.NoError(setErr)
	setErr = testProfile.SetAttribute(global.TEST_VAR_FIRST_SECTION_KEY, global.TEST_VAR_FIRST_SECTION_UNIQUE_KEY_LABEL, global.TEST_VAR_FIRST_SECTION_UNIQUE_KEY_VALUE)
	assert.NoError(setErr)
	setErr = testProfile.SetAttribute(global.TEST_VAR_SECOND_SECTION_KEY, global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.NoError(setErr)
	setErr = testProfile.SetAttribute(global.TEST_VAR_SECOND_SECTION_KEY, global.TEST_VAR_SECOND_SECTION_UNIQUE_KEY_LABEL, global.TEST_VAR_SECOND_SECTION_UNIQUE_KEY_VALUE)
	assert.NoError(setErr)

	allAttrs := testProfile.GetAllAttributes()
	assert.Equal(allAttrs[global.SECTION_NAME_BLANK][global.TEST_VAR_ATTRIBUTE_NAME_LABEL], global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.Equal(allAttrs[global.TEST_VAR_FIRST_SECTION_KEY][global.TEST_VAR_ATTRIBUTE_NAME_LABEL], global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.Equal(allAttrs[global.TEST_VAR_FIRST_SECTION_KEY][global.TEST_VAR_FIRST_SECTION_UNIQUE_KEY_LABEL], global.TEST_VAR_FIRST_SECTION_UNIQUE_KEY_VALUE)
	assert.Equal(allAttrs[global.TEST_VAR_SECOND_SECTION_KEY][global.TEST_VAR_ATTRIBUTE_NAME_LABEL], global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.Equal(allAttrs[global.TEST_VAR_SECOND_SECTION_KEY][global.TEST_VAR_SECOND_SECTION_UNIQUE_KEY_LABEL], global.TEST_VAR_SECOND_SECTION_UNIQUE_KEY_VALUE)

	os.RemoveAll(testFactory.ParentDirectory)
}

func TestProfileSectionAttributes(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing getting section attributes of a profile.")
	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(factoryErr)

	testProfile, newErr := New(global.TEST_VAR_FIRST_PROFILE_LABEL, testFactory)
	assert.NoError(newErr)
	assert.Equal(global.TEST_VAR_FIRST_PROFILE_LABEL, testProfile.Name)
	assert.Equal(testFactory.ConfigDirectory+global.TEST_VAR_FIRST_PROFILE_LABEL, testProfile.ConfigFileLocation)
	assert.Equal(testProfile.Factory, testFactory)
	assert.True(testProfile.Initialized)

	setErr := testProfile.SetAttribute(global.TEST_VAR_FIRST_SECTION_KEY, global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.NoError(setErr)
	setErr = testProfile.SetAttribute(global.TEST_VAR_FIRST_SECTION_KEY, global.TEST_VAR_FIRST_SECTION_UNIQUE_KEY_LABEL, global.TEST_VAR_FIRST_SECTION_UNIQUE_KEY_VALUE)
	assert.NoError(setErr)

	sectionAttrs, sectErr := testProfile.GetAllSectionAttributes(global.TEST_VAR_FIRST_SECTION_KEY)
	assert.NoError(sectErr)
	assert.Equal(sectionAttrs[global.TEST_VAR_ATTRIBUTE_NAME_LABEL], global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.Equal(sectionAttrs[global.TEST_VAR_FIRST_SECTION_UNIQUE_KEY_LABEL], global.TEST_VAR_FIRST_SECTION_UNIQUE_KEY_VALUE)

	os.RemoveAll(testFactory.ParentDirectory)
}
