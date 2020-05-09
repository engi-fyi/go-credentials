package profile

import (
	"errors"
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"github.com/rs/zerolog/log"
	"os"
	"regexp"
)

func New(profileName string, credentialFactory *factory.Factory) (*Profile, error) {
	keyRegex := regexp.MustCompile(global.REGEX_KEY_NAME)

	if !keyRegex.MatchString(profileName) {
		return nil, errors.New(ERR_PROFILE_NAME_MUST_MATCH_REGEX)
	}

	newProfile := Profile{
		Name: profileName,
		ConfigFileLocation: credentialFactory.ConfigDirectory + profileName,
		attributes: make(map[string]map[string]string),
		Initialized: true,
		Factory: credentialFactory,
	}

	return &newProfile, nil
}

func Remove(thisProfile *Profile) error {
	log.Trace().Str("profile", thisProfile.Name).Msg("Deleting profile.")
	removeErr := os.Remove(thisProfile.ConfigFileLocation)
	*thisProfile = Profile{}
	return removeErr
}

func (thisProfile *Profile) SetAttribute(sectionName string, key string, value string) error {
	keyRegex := regexp.MustCompile(global.REGEX_KEY_NAME)

	if sectionName == "" {
		sectionName = NO_SECTION_KEY
	}

	if !keyRegex.MatchString(sectionName) || !keyRegex.MatchString(key) {
		return errors.New(ERR_MUST_MATCH_REGEX)
	}

	if _, ok := thisProfile.attributes[sectionName]; !ok {
		thisProfile.attributes[sectionName] = make (map[string]string)
	}

	thisProfile.attributes[sectionName][key] = value
	return nil
}

func (thisProfile *Profile) GetAttribute(sectionName string, key string) string {
	if sectionName == "" {
		sectionName = NO_SECTION_KEY
	}

	if _, ok := thisProfile.attributes[sectionName]; ok {
		if _, ok := thisProfile.attributes[sectionName][key]; ok {
			return thisProfile.attributes[sectionName][key]
		}
	}

	return "" // attribute cannot have a blank value, so always assume a blank value was a NOT_FOUND err
}

func (thisProfile *Profile) DeleteAttribute(sectionName string, key string) error {
	if len(thisProfile.GetAttribute(sectionName, key)) == 0 {
		return errors.New(ERR_DELETED_ATTRIBUTE_NOT_EXIST)
	}

	delete(thisProfile.attributes[sectionName], key)
	return nil
}