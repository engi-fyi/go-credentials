package credential

import (
	"errors"
	"fmt"
	"github.com/engi-fyi/go-credentials/environment"
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"github.com/rs/zerolog/log"
	"gopkg.in/ini.v1"
	"os"
	"regexp"
	"strings"
	"time"
)

// New creates a new Credential instance. credentialFactory provides global application-level settings for the
// go-credential library. username and password are the user's base credentials. No other attributes are set during
// the initial creation of a Credential object. It returns a pointer to a new Credential object.
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

	return &Credential{
		Username:    username,
		Password:    password,
		Initialized: true,
		Factory:     credentialFactory,
		attributes:  make(map[string]string),
	}, nil
}

// SetAttribute sets an attribute on a Credential object. key must match the regex '(?m)^[0-9A-Za-z_]+$'. There are no
// restrictions on the value of an attribute, aside from Go-level restrictions on strings. When these values are
// processed by Save(), they are stored under the [attributes] category. Attributes set by SetAttribute can only be
// accessed by the sister function GetAttribute. If username or password is passed as the attribute key, the set is
// redirected to the Username or Password property on the Credential object. This returns an error (if applicable).
func (credential *Credential) SetAttribute(key string, value string) error {
	if strings.ToLower(key) == global.USERNAME_LABEL {
		log.Trace().Msg("Redirected attribute request to set username.")
		credential.Username = value
	} else if strings.ToLower(key) == global.PASSWORD_LABEL {
		log.Trace().Msg("Redirected attribute request to set password.")
		credential.Password = value
	} else {
		keyRegex := regexp.MustCompile(global.REGEX_KEY_NAME)

		if keyRegex.MatchString(key) {
			log.Trace().Str("key", key).Msg("Setting attribute.")
			credential.attributes[key] = value
		} else {
			return errors.New(ERR_KEY_MUST_MATCH_REGEX)
		}
	}

	return nil
}

// GetAttribute retrieves an attribute that has been stored in the unexported Credential property attributes. A key
// is passed in, and if the key does not have a value stored in the Credential, an error is returned.
func (credential *Credential) GetAttribute(key string) (string, error) {
	if strings.ToLower(key) == global.USERNAME_LABEL {
		log.Trace().Msg("Redirected attribute request to value of username.")
		return credential.Username, nil
	} else if strings.ToLower(key) == global.PASSWORD_LABEL {
		log.Trace().Msg("Redirected attribute request to value of password.")
		return credential.Password, nil
	} else {
		log.Trace().Str("Attribute", key).Msg("Retrieving attribute.")
		value, exists := credential.attributes[key]

		if exists {
			return value, nil
		} else {
			log.Trace().Str("Attribute", key).Msg("That attribute doesn't exist.")
			return "", errors.New(ERR_ATTRIBUTE_NOT_EXIST)
		}
	}
}

// LoadFromEnvironment is responsible for scanning environment variables and retrieves applicable variables that have
// the prefix of the application name that has been set in the credentialFactory. Most importantly, it will scan for
// the keys username or password; if an alternate for either of these have been set, those will be loaded instead.
// The respective properties Username and Password on the Credential object will be set. The rest of the variables will
// be stored as attributes and be accessible via GetAttribute.
//
// Example 1: Normal Usage
//
// If credentialFactory.ApplicationName has been set to TEST_APP, any environment variables beginning with
// TEST_APP will be imported into the credential object.
//
// Example 2: Alternates Usage
//
// If credentialFactory.Alternates["username"] has been set to ACCESS_TOKEN, then if an environment variable named
// TEST_APP_ACCESS_TOKEN exists its value will be stored in the resulting Credential object's Username property.
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

// Save is responsible for saving the credential at ~/.application_name/credentials in the specified output format
// that has been set on the Credentials' Factory object.
//
// TODO(7): Implement json format.
func (credential *Credential) Save() error {
	if credential.Factory.Initialized {
		if credential.Factory.OutputType == global.OUTPUT_TYPE_JSON {
			return credential.saveJson()
		} else if credential.Factory.OutputType == global.OUTPUT_TYPE_INI {
			return credential.saveIni()
		} else {
			return errors.New(factory.ERR_FACTORY_INCONSISTENT_STATE)
		}
	} else {
		return errors.New(factory.ERR_FACTORY_NOT_INITIALIZED)
	}
}

// Load is the default method of loading existing credentials. First it will attempt to load credentials from the
// environment, as they take precedence over file-based credentials. If that is not successful, it will attempt to
// load credentials from the credentials file using the appropriate format set in the sourceFactory. If the file doesn't
// exist or is in an unexpected format, an error will be thrown.
//
// TODO(7): Implement json format.
func Load(sourceFactory *factory.Factory) (*Credential, error) {
	var fileErr error
	loadedCredential, envErr := LoadFromEnvironment(sourceFactory)

	if envErr != nil || !loadedCredential.Initialized {
		if sourceFactory.OutputType == global.OUTPUT_TYPE_JSON {
			loadedCredential, fileErr = loadJson(sourceFactory)
		} else if sourceFactory.OutputType == global.OUTPUT_TYPE_INI {
			loadedCredential, fileErr = LoadFromIniFile(sourceFactory)
		} else {
			return nil, errors.New(factory.ERR_FACTORY_INCONSISTENT_STATE)
		}
	}

	if envErr != nil && fileErr != nil {
		return nil, fmt.Errorf("%v. %v", envErr, fileErr)
	}

	return loadedCredential, nil
}

// DeployEnv is used in cases where your application requires environment variables, but your users configure their
// credentials via file-based methods. Simply Load the Credential from a file, then DeployEnv will deploy the variables
// as follows:
// - Username: APP_NAME_USERNAME
// - Other Attributes: APP_NAME_ATTRIBUTE_NAME
//
// Currently, we do not export the Password property to the environment, but it is in the pipeline to enable this in
// the future.
//
// TODO(3): Export Password via Boolean
func (credential *Credential) DeployEnv() error {
	if credential.Factory.UseEnvironment {
		setKeys, deployErr := environment.Deploy(
			credential.Factory.ApplicationName,
			credential.Username,
			"",
			credential.Factory.GetAlternates(),
			credential.attributes)

		credential.environmentVariables = setKeys
		return deployErr
	} else {
		return errors.New(ERR_FACTORY_PRIVATE_ATTEMPT_DEPLOY)
	}
}

// GetEnvironmentVariables retrieves a list of internally managed environment variables that have been set by
// go-credentials. This value only has a value if DeployEnv has been used.
func (credential *Credential) GetEnvironmentVariables() []string {
	return credential.environmentVariables
}

// BUG(4): Respect alternates when saving to file.
func (credential *Credential) saveIni() error {
	if credential.Factory.Initialized {
		log.Trace().Msg("Saving credentials as ini.")
		credentialFile, loadErr := loadIni(credential.Factory)

		if loadErr != nil {
			return loadErr
		}

		log.Trace().Msg("Existing credentials loaded from file.")

		credentialFile.Section("default").Key("username").SetValue(credential.Username)
		credentialFile.Section("default").Key("password").SetValue(credential.Password)
		log.Trace().Msg("Username and password set in ini.")

		for attributeKey, attributeValue := range credential.attributes {
			credentialFile.Section("attributes").Key(attributeKey).SetValue(attributeValue)
			log.Trace().Str("Attribute", attributeKey).Msg("Setting ini attribute.")
		}

		log.Trace().Msg("Setting ini file metadata.")
		credentialFile.Section("metadata").Key("last_updated").SetValue(time.Now().String())
		saveErr := credentialFile.SaveTo(credential.Factory.CredentialFile)

		if saveErr != nil {
			log.Error().Err(saveErr).Msg("Error saving credentials ini file.")
			return saveErr
		}

		log.Trace().Msg("Credential ini file saved successfully.")

		return nil
	} else {
		return errors.New(factory.ERR_FACTORY_NOT_INITIALIZED)
	}
}

// LoadFromIniFile is responsible for saving a Credential object as in the ini format at ~/.app_name/credentials.
// username and password are stored under the default heading, and all attributes are stored under the attributes
// heading
//
// Example Credential File
// 		[default]
//		username=my_username
//		password=my_password01
//
//		[attributes]
//		an_attribute=the value of an attribute
//
// BUG(4): Respect alternates when loading from file.
func LoadFromIniFile(fromFactory *factory.Factory) (*Credential, error) {
	if fromFactory.Initialized {
		credentialFile, loadErr := loadIni(fromFactory)

		if loadErr != nil {
			log.Error().Err(loadErr).Msg("Error loading credentials.")

			return nil, loadErr
		}

		loadedCredential, credErr := getCredentialFromIni(fromFactory, credentialFile)

		if credErr != nil {
			return nil, credErr
		}

		attrErr := addAttributesFromIni(credentialFile, loadedCredential)

		if attrErr != nil {
			return nil, attrErr
		}

		return loadedCredential, nil
	} else {
		return nil, errors.New(factory.ERR_FACTORY_NOT_INITIALIZED)
	}
}

func loadIni(fromFactory *factory.Factory) (*ini.File, error) {
	credentialFile, loadErr := ini.Load(fromFactory.CredentialFile)

	if _, cfsErr := os.Stat(fromFactory.CredentialFile); os.IsNotExist(cfsErr) {
		log.Trace().Msg("File doesn't exist, creating then existing.")
		if _, cfsErr := os.Stat(fromFactory.CredentialFile); os.IsNotExist(cfsErr) {
			log.Trace().Str("Credential File", fromFactory.CredentialFile).Msg("Initializing credential file.")
			emptyFile, efErr := os.Create(fromFactory.CredentialFile)

			if efErr != nil {
				log.Error().Err(efErr).Msg("Error creating credential file.")
				return nil, efErr
			}

			closeErr := emptyFile.Close()

			if closeErr != nil {
				return nil, closeErr
			}

			credentialFile, loadErr = ini.Load(fromFactory.CredentialFile)

			if loadErr != nil {
				return nil, loadErr
			}
		}
	}

	if loadErr != nil {
		return nil, loadErr
	} else {
		return credentialFile, nil
	}
}

func getCredentialFromIni(fromFactory *factory.Factory, credentialFile *ini.File) (*Credential, error) {
	log.Trace().Msg("Creating credential from ini.")
	loadedCredential, newErr := New(fromFactory,
		credentialFile.Section("default").Key("username").String(),
		credentialFile.Section("default").Key("password").String())

	if newErr != nil {
		return nil, newErr
	}

	return loadedCredential, nil
}

func addAttributesFromIni(credentialFile *ini.File, loadedCredential *Credential) error {
	attributes, attributeErr := credentialFile.GetSection("attributes")

	if attributeErr != nil {
		if attributeErr.Error() != "section \"attributes\" does not exist" {
			return attributeErr
		} else {
			log.Trace().Msg("attributes section doesn't exist, skipping.")
		}
	} else {
		attributeKeys := attributes.KeyStrings()

		for i := 0; i < len(attributeKeys); i++ {
			attributeKey := attributeKeys[i]
			attributeValue := credentialFile.Section("attributes").Key(attributeKey).String()
			log.Trace().Str("Attribute Key", attributeKey).Msg("Loading value for attribute.")
			loadedCredential.attributes[attributeKey] = attributeValue
		}
	}

	return nil
}

// TODO(7): Implement json format.
func (credential *Credential) saveJson() error {
	return errors.New(ERR_JSON_FUNCTIONALITY_NOT_IMPLEMENTED)
}

// TODO(7): Implement json format.
func loadJson(factory *factory.Factory) (*Credential, error) {
	return nil, errors.New(ERR_JSON_FUNCTIONALITY_NOT_IMPLEMENTED)
}
