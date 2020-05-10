package credential

import (
	"errors"
	"regexp"
	"strings"

	"github.com/engi-fyi/go-credentials/environment"
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"github.com/engi-fyi/go-credentials/profile"
	"github.com/rs/zerolog/log"
)

/*
New creates a new Credential instance. credentialFactory provides global application-level settings for the
go-credential library. username and password are the user's base credentials. No other attributes are set during
the initial creation of a Credential object, and these are stored in the Profile.
*/
func New(credentialFactory *factory.Factory, username string, password string) (*Credential, error) {
	log.Trace().Msg("Building credential object.")

	if credentialFactory == nil {
		return nil, errors.New(ERR_FACTORY_MUST_BE_INITIALIZED)
	} else {
		if !credentialFactory.Initialized {
			return nil, errors.New(ERR_FACTORY_MUST_BE_INITIALIZED)
		}
	}

	if username == "" || password == "" {
		return nil, errors.New(ERR_USERNAME_OR_PASSWORD_NOT_SET)
	}

	newCredential := &Credential{
		Username:    username,
		Password:    password,
		Initialized: true,
		Factory:     credentialFactory,
	}
	profileErr := newCredential.SetProfile(global.DEFAULT_PROFILE_NAME)

	if profileErr != nil {
		return nil, profileErr
	}

	return newCredential, nil
}

/*
SetProfile sets the active profile on the Credential object. It does not change the value of the Credential's username
or password, and only manages the attached Profile object.
*/
func (thisCredential *Credential) SetProfile(profileName string) error {
	if !thisCredential.Initialized {
		return errors.New(ERR_CREDENTIAL_NOT_INITIALIZED)
	}

	log.Trace().Msg("Attempting to load profile from file.")
	myProfile, profileErr := profile.Load(profileName, thisCredential.Factory)

	if profileErr != nil {
		if profileErr.Error() == profile.ERR_PROFILE_DID_NOT_EXIST {
			log.Trace().Msg("Profile did not exist, creating new.")
			myProfile, profileErr = profile.New(profileName, thisCredential.Factory)
		} else {
			log.Error().Err(profileErr).Msg("Sorry there was an error loading the profile.")
			return profileErr
		}
	}

	if !myProfile.Initialized {
		log.Error().Err(profileErr).Msg("The profile has not been initialized.")
		return errors.New(profile.ERR_PROFILE_NOT_INITIALIZED)
	}

	thisCredential.Profile = myProfile
	return nil
}

/*
SetAttribute sets an attribute on a Credential object. key must match the regex '(?m)^[0-9A-Za-z_]+$'. There are no
restrictions on the value of an attribute, aside from Go-level restrictions on strings. When these values are processed
by Save(), they are stored in the config file for the Credential's currently set Profile. If username or password is
passed as the attribute key, the set is redirected to the Username or Password property on the Credential object.
*/
func (thisCredential *Credential) SetAttribute(key string, value string) error {
	if strings.ToLower(key) == global.USERNAME_LABEL {
		log.Trace().Msg("Redirected attribute request to set username.")
		thisCredential.Username = value
	} else if strings.ToLower(key) == global.PASSWORD_LABEL {
		log.Trace().Msg("Redirected attribute request to set password.")
		thisCredential.Password = value
	} else {
		keyRegex := regexp.MustCompile(global.REGEX_KEY_NAME)

		if keyRegex.MatchString(key) {
			log.Trace().Str("key", key).Msg("Setting attribute.")

			attributeErr := thisCredential.Profile.SetAttribute(global.NO_SECTION_KEY, key, value)

			if attributeErr != nil {
				return attributeErr
			}
		} else {
			return errors.New(ERR_KEY_MUST_MATCH_REGEX)
		}
	}

	return nil
}

/*
GetAttribute retrieves an attribute that has been stored on the Credential's associated Profile. A key
is passed in, and if the key does not have a value stored in the Profile, an error is returned.
*/
func (thisCredential *Credential) GetAttribute(key string) (string, error) {
	if strings.ToLower(key) == global.USERNAME_LABEL {
		log.Trace().Msg("Redirected attribute request to value of username.")
		return thisCredential.Username, nil
	} else if strings.ToLower(key) == global.PASSWORD_LABEL {
		log.Trace().Msg("Redirected attribute request to value of password.")
		return thisCredential.Password, nil
	} else {
		log.Trace().Str("Attribute", key).Msg("Retrieving attribute.")
		value := thisCredential.Profile.GetAttribute(global.NO_SECTION_KEY, key)

		if len(value) > 0 {
			return value, nil
		} else {
			log.Trace().Str("Attribute", key).Msg("That attribute doesn't exist.")
			return "", errors.New(ERR_ATTRIBUTE_NOT_EXIST)
		}
	}
}

/*
LoadFromEnvironment is responsible for scanning environment variables and retrieves applicable variables that have
the prefix of the application name that has been set in the credentialFactory. Most importantly, it will scan for
the keys username or password; if an alternate for either of these have been set, those will be loaded instead.
The respective properties Username and Password on the Credential object will be set. The rest of the variables will
be stored as attributes and be accessible via GetAttribute. These are loaded into the Profile object and can also
be accessed by its relevant functions.

Example 1: Normal Usage

If credentialFactory.ApplicationName has been set to TEST_APP, any environment variables beginning with
TEST_APP will be imported into the Credential (and subsequently Profile) object.

Example 2: Alternates Usage

If credentialFactory.Alternates["username"] has been set to ACCESS_TOKEN, then if an environment variable named
TEST_APP_ACCESS_TOKEN exists its value will be stored in the resulting Credential object's Username property.
*/
func LoadFromEnvironment(credentialFactory *factory.Factory) (*Credential, error) {
	log.Trace().Msg("Creating credentials from environment variable.")
	values, loadErr := environment.Load(credentialFactory.ApplicationName, credentialFactory.GetAlternates())
	username := ""
	password := ""

	if loadErr != nil {
		return nil, loadErr
	}

	for key, value := range values {
		if key == global.USERNAME_LABEL || key == credentialFactory.GetAlternateUsername() {
			log.Trace().Msg("Found username key.")
			username = value
		}

		if key == global.PASSWORD_LABEL || key == credentialFactory.GetAlternatePassword() {
			log.Trace().Msg("Found password key.")
			password = value
		}
	}

	credential, credErr := New(credentialFactory, username, password)

	if credErr != nil {
		return nil, credErr
	}

	profileErr := credential.SetProfile("ENVIRONMENT")

	if profileErr != nil {
		return nil, profileErr
	}

	log.Trace().Msg("Adding attributes.")
	for key, value := range values {
		if key != global.USERNAME_LABEL && key != global.PASSWORD_LABEL {
			log.Trace().Str("key", key).Msg("Found attribute, setting.")
			attrErr := credential.SetAttribute(key, value)

			if attrErr != nil {
				return nil, attrErr
			}
		}
	}

	return credential, nil
}

/*DeployEnv is used in cases where your application requires environment variables, but your users configure their
credentials via file-based methods. Simply Load the Credential from a file, then DeployEnv will deploy the variables
as follows:
	- Username: APP_NAME_USERNAME
	- Other Attributes: APP_NAME_ATTRIBUTE_NAME

Currently, we do not export the Password property to the environment, but it is in the pipeline to enable this in
the future.
*/
// TODO(3): Export Password via Boolean
func (thisCredential *Credential) DeployEnv() error {
	if thisCredential.Factory.UseEnvironment {
		setKeys, deployErr := environment.Deploy(
			thisCredential.Factory.ApplicationName,
			thisCredential.Username,
			"",
			thisCredential.Factory.GetAlternates(),
			thisCredential.Profile.GetAllAttributes())

		thisCredential.environmentVariables = setKeys
		return deployErr
	} else {
		return errors.New(ERR_FACTORY_PRIVATE_ATTEMPT_DEPLOY)
	}
}

/*
GetEnvironmentVariables retrieves a list of internally managed environment variables that have been set by
go-credentials. This value only has a value if DeployEnv has been used.
*/
func (thisCredential *Credential) GetEnvironmentVariables() []string {
	return thisCredential.environmentVariables
}

// TODO(7): Implement json format.
func loadJson(factory *factory.Factory) (*Credential, error) {
	return nil, errors.New(ERR_JSON_FUNCTIONALITY_NOT_IMPLEMENTED)
}
