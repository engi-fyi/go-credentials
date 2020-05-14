package credential

import (
	"errors"
	"regexp"
	"strings"

	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"github.com/engi-fyi/go-credentials/profile"
)

/*
New creates a new Credential instance. credentialFactory provides global application-level settings for the
go-credential library. username and password are the user's base credentials. No other attributes are set during
the initial creation of a Credential object, and these are stored in the Profile.
*/

func New(sourceFactory *factory.Factory, username string, password string) (*Credential, error) {
	return NewProfile(global.DEFAULT_PROFILE_NAME, sourceFactory, username, password)
}

func NewProfile(profileName string, sourceFactory *factory.Factory, username string, password string) (*Credential, error) {
	if sourceFactory == nil {
		return nil, errors.New(ERR_FACTORY_MUST_BE_INITIALIZED)
	} else {
		if !sourceFactory.Initialized {
			return nil, errors.New(ERR_FACTORY_MUST_BE_INITIALIZED)
		}
	}

	sourceFactory.Log.Trace().Str("profile", profileName).Msg("Loading credentials from default profile")
	sourceFactory.Log.Trace().Msg("Building credential object.")

	if username == "" || password == "" {
		return nil, errors.New(ERR_USERNAME_OR_PASSWORD_NOT_SET)
	}

	newProfile, profileErr := profile.New(profileName, sourceFactory)

	if profileErr != nil {
		return nil, profileErr
	}

	newCredential := &Credential{
		Username:    username,
		Password:    password,
		Initialized: true,
		Factory:     sourceFactory,
		Profile:     newProfile,
	}

	return newCredential, nil
}

/*
SetAttribute sets an attribute on a Credential object. key must match the regex '(?m)^[0-9A-Za-z_]+$'. There are no
restrictions on the value of an attribute, aside from Go-level restrictions on strings. When these values are processed
by Save(), they are stored in the config file for the Credential's currently set Profile. If username or password is
passed as the attribute key, the set is redirected to the Username or Password property on the Credential object.
*/
func (thisCredential *Credential) SetAttribute(key string, value string) error {
	if strings.ToLower(key) == global.USERNAME_LABEL {
		if thisCredential.selectedSection == "" { // can only get here if .Section() is not used.
			thisCredential.Factory.Log.Trace().Msg("Redirected attribute request to set username.")
			thisCredential.Username = value
		} else {
			thisCredential.Factory.Log.Error().Msg("Cannot redirect attribute request to username because a section has been specified.")
			return errors.New(ERR_CANNOT_SET_USERNAME_WHEN_USING_SECTION)
		}

	} else if strings.ToLower(key) == global.PASSWORD_LABEL {
		if thisCredential.selectedSection == "" { // can only get here if .Section() is not used.
			thisCredential.Factory.Log.Trace().Msg("Redirected attribute request to set password.")
			thisCredential.Password = value
		} else {
			thisCredential.Factory.Log.Error().Msg("Cannot redirect attribute request to password because a section has been specified.")
			return errors.New(ERR_CANNOT_SET_PASSWORD_WHEN_USING_SECTION)
		}
	} else {
		keyRegex := regexp.MustCompile(global.REGEX_KEY_NAME)

		if keyRegex.MatchString(key) {
			thisCredential.Factory.Log.Trace().Str("key", key).Msg("Setting attribute.")
			var attributeErr error

			if thisCredential.selectedSection == global.SECTION_NAME_BLANK {
				attributeErr = thisCredential.Profile.SetAttribute(global.NO_SECTION_KEY, key, value)
			} else {
				attributeErr = thisCredential.Profile.SetAttribute(thisCredential.selectedSection, key, value)
			}

			if attributeErr != nil {
				thisCredential.Factory.Log.Error().Err(attributeErr).Str("key", key).Msg("Error setting attribute.")
				return attributeErr
			}
		} else {
			thisCredential.Factory.Log.Error().Msg(ERR_KEY_MUST_MATCH_REGEX)
			return errors.New(ERR_KEY_MUST_MATCH_REGEX)
		}
	}

	return nil
}

/*
Section is used as an intermediary function between a credential and GetAttribute or SetAttribute. When you pass a
section name, the subsequent method attribute function will be called against that section.
*/
func (thisCredential *Credential) Section(section string) *Credential {
	credentialWithSectionSet := *thisCredential

	if section == "" {
		credentialWithSectionSet.selectedSection = global.SECTION_NAME_BLANK
	} else {
		credentialWithSectionSet.selectedSection = section
	}

	return &credentialWithSectionSet
}

/*
GetAttribute retrieves an attribute that has been stored on the Credential's associated Profile. A key
is passed in, and if the key does not have a value stored in the Profile, an error is returned.
*/
func (thisCredential *Credential) GetAttribute(key string) string {
	if strings.ToLower(key) == global.USERNAME_LABEL {
		thisCredential.Factory.Log.Trace().Msg("Redirected attribute request to value of username.")
		return thisCredential.Username
	} else if strings.ToLower(key) == global.PASSWORD_LABEL {
		thisCredential.Factory.Log.Trace().Msg("Redirected attribute request to value of password.")
		return thisCredential.Password
	} else {
		thisCredential.Factory.Log.Trace().Str("key", key).Msg("Retrieving attribute.")
		var value string

		if thisCredential.selectedSection == "" {
			value = thisCredential.Profile.GetAttribute(global.NO_SECTION_KEY, key)
		} else {
			value = thisCredential.Profile.GetAttribute(thisCredential.selectedSection, key)
		}

		if len(value) > 0 {
			return value
		} else {
			thisCredential.Factory.Log.Trace().Str("key", key).Msg("That attribute doesn't exist.")
			return ""
		}
	}
}

/*
DeleteAttribute removes an attribute that has been stored on the Credential's associated Profile. If the value is
deleted an error is not returned. Otherwise, an error will be returned.
*/
func (thisCredential *Credential) DeleteAttribute(key string) error {
	if strings.ToLower(key) == global.USERNAME_LABEL {
		thisCredential.Factory.Log.Trace().Msg("Can't delete username from a Credential.")
		return errors.New(ERR_CANNOT_REMOVE_USERNAME)
	} else if strings.ToLower(key) == global.PASSWORD_LABEL {
		thisCredential.Factory.Log.Trace().Msg("Can't delete password from a Credential.")
		return errors.New(ERR_CANNOT_REMOVE_PASSWORD)
	} else {
		thisCredential.Factory.Log.Trace().Str("key", key).Msg("Deleting attribute.")

		if thisCredential.selectedSection == "" {
			return thisCredential.Profile.DeleteAttribute(global.NO_SECTION_KEY, key)
		} else {
			return thisCredential.Profile.DeleteAttribute(thisCredential.selectedSection, key)
		}
	}
}
