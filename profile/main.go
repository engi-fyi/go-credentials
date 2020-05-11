package profile

import (
	"errors"
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"regexp"
)

/*
New is responsible for constructing a new, blank profile to be used by a Credential. It is important to note, that
this function does not save a profile, and this needs to be done using the Save() function.
*/
func New(profileName string, credentialFactory *factory.Factory) (*Profile, error) {
	keyRegex := regexp.MustCompile(global.REGEX_KEY_NAME)

	if !keyRegex.MatchString(profileName) {
		return nil, errors.New(ERR_PROFILE_NAME_MUST_MATCH_REGEX)
	}

	newProfile := Profile{
		Name:               profileName,
		ConfigFileLocation: credentialFactory.ConfigDirectory + profileName,
		attributes:         make(map[string]map[string]string),
		Initialized:        true,
		Factory:            credentialFactory,
	}

	return &newProfile, nil
}

/*
SetAttribute is responsible for setting an attribute against a section name. If the section name is blank, the attribute
will be stored without a section.

Example: With Section
	myProfile.SetAttribute("a_section", "a_key", "a_value")

	// my_profile (ini)
	// [a_section]
	// a_key = a_value <-- this value will be set

Example: No Section
	myProfile.SetAttribute("", "a_key", "a_value")

	// my_profile (ini)
	// a_key = a_value <-- this value will be set

	// [a_section]
	// a_key = a_value
*/
func (thisProfile *Profile) SetAttribute(sectionName string, key string, value string) error {
	keyRegex := regexp.MustCompile(global.REGEX_KEY_NAME)

	if sectionName == "" {
		sectionName = global.NO_SECTION_KEY
	}

	if !keyRegex.MatchString(sectionName) || !keyRegex.MatchString(key) {
		return errors.New(ERR_MUST_MATCH_REGEX)
	}

	if _, ok := thisProfile.attributes[sectionName]; !ok {
		thisProfile.attributes[sectionName] = make(map[string]string)
	}

	thisProfile.attributes[sectionName][key] = value
	return nil
}

/*
GetAttribute retrieves an attribute from a profile section. If the section name is blank, the attribute will be
retrieved from the default store which has no section name.
*/
func (thisProfile *Profile) GetAttribute(sectionName string, key string) string {
	if sectionName == "" {
		sectionName = global.NO_SECTION_KEY
	}

	if _, ok := thisProfile.attributes[sectionName]; ok {
		if _, ok := thisProfile.attributes[sectionName][key]; ok {
			return thisProfile.attributes[sectionName][key]
		}
	}

	return "" // attribute cannot have a blank value, so always assume a blank value was a NOT_FOUND err
}

/*
GetAllAttribute simply returns the a nested map that has two sets of keys, a section key, and then a normal key. Any
attributes without a section will be returned under the section key of "default".

Example: Map Structure
	default:
		my_key: my_value
	a_section:
		my_key: my_value
		a_key: a_value
	b_section:
		my_key: my_value
		b_key: b_value
*/
func (thisProfile *Profile) GetAllAttributes() map[string]map[string]string {
	return thisProfile.attributes
}

/*
GetAllSection attributes is responsible for returning all of the attributes for one section. It returns this as a simple
map with the key and value for each attribute only. There are no section references returned by this function.
*/
func (thisProfile *Profile) GetAllSectionAttributes(sectionName string) (map[string]string, error) {
	if _, ok := thisProfile.attributes[sectionName]; ok {
		return thisProfile.attributes[sectionName], nil
	} else {
		return nil, errors.New(ERR_SECTION_NOT_EXIST)
	}
}

/*
DeleteAttribute removes an attribute from a profile. It is important to note that until the Profile is Save()d, the
attribute may still exist on the file system.
*/
func (thisProfile *Profile) DeleteAttribute(sectionName string, key string) error {
	if len(thisProfile.GetAttribute(sectionName, key)) == 0 {
		return errors.New(ERR_DELETED_ATTRIBUTE_NOT_EXIST)
	}

	delete(thisProfile.attributes[sectionName], key)
	return nil
}
