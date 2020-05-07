package credential

import (
	"errors"
	"fmt"
	"github.com/HammoTime/go-credentials/environment"
	"github.com/HammoTime/go-credentials/factory"
	"github.com/rs/zerolog/log"
	"gopkg.in/ini.v1"
	"os"
	"regexp"
	"strings"
	"time"
)

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

func (credential *Credential) SetAttribute(key string, value string) error {
	if strings.ToLower(key) == ATTRIBUTE_USERNAME {
		log.Trace().Msg("Redirected attribute request to set username.")
		credential.Username = value
	} else if strings.ToLower(key) == ATTRIBUTE_PASSWORD {
		log.Trace().Msg("Redirected attribute request to set password.")
		credential.Password = value
	} else {
		keyRegex := regexp.MustCompile(REGEX_KEY_NAME)

		if keyRegex.MatchString(key) {
			log.Trace().Str("key", key).Msg("Setting attribute.")
			credential.attributes[key] = value
		} else {
			return errors.New(ERR_KEY_MUST_MATCH_REGEX)
		}
	}

	return nil
}

func (credential *Credential) GetAttribute(key string) (string, error) {
	if strings.ToLower(key) == ATTRIBUTE_USERNAME {
		log.Trace().Msg("Redirected attribute request to value of username.")
		return credential.Username, nil
	} else if strings.ToLower(key) == ATTRIBUTE_PASSWORD {
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

func LoadFromEnvironment(credentialFactory *factory.Factory) (*Credential, error) {
	log.Trace().Msg("Creating credentials from environment variable.")
	values, loadErr := environment.Load(credentialFactory.ApplicationName, credentialFactory.Alternates)
	username := ""
	password := ""

	if loadErr != nil {
		return nil, loadErr
	}

	for key, value := range values {
		if key == ATTRIBUTE_USERNAME || key == credentialFactory.Alternates[ATTRIBUTE_USERNAME] {
			log.Trace().Msg("Found username key.")
			username = value
		}

		if key == ATTRIBUTE_PASSWORD || key == credentialFactory.Alternates[ATTRIBUTE_PASSWORD] {
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
		if key != ATTRIBUTE_USERNAME && key != ATTRIBUTE_PASSWORD {
			log.Trace().Str("key", key).Msg("Found attribute, setting.")
			attrErr := credential.SetAttribute(key, value)

			if attrErr != nil {
				return nil, attrErr
			}
		}
	}

	return credential, nil
}

func (credential *Credential) Save() error {
	if credential.Factory.Initialized {
		if credential.Factory.OutputType == factory.OUTPUT_TYPE_JSON {
			return credential.saveJson()
		} else if credential.Factory.OutputType == factory.OUTPUT_TYPE_INI {
			return credential.saveIni()
		} else {
			return errors.New(factory.ERR_FACTORY_INCONSISTENT_STATE)
		}
	} else {
		return errors.New(factory.ERR_FACTORY_NOT_INITIALIZED)
	}
}

func Load(sourceFactory *factory.Factory) (*Credential, error) {
	var fileErr error
	loadedCredential, envErr := LoadFromEnvironment(sourceFactory)

	if envErr != nil || !loadedCredential.Initialized {
		if sourceFactory.OutputType == factory.OUTPUT_TYPE_JSON {
			loadedCredential, fileErr = loadJson(sourceFactory)
		} else if sourceFactory.OutputType == factory.OUTPUT_TYPE_INI {
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

func (credential *Credential) DeployEnv() error {
	if credential.Factory.UseEnvironment {
		setKeys, deployErr := environment.Deploy(
			credential.Factory.ApplicationName,
			credential.Username,
			"",
			credential.Factory.Alternates,
			credential.attributes)

		credential.environmentVariables = setKeys
		return deployErr
	} else {
		return errors.New(ERR_FACTORY_PRIVATE_ATTEMPT_DEPLOY)
	}
}

func (credential *Credential) GetEnvironmentVariables() []string {
	return credential.environmentVariables
}

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

func (credential *Credential) saveJson() error {
	return errors.New(ERR_JSON_FUNCTIONALITY_NOT_IMPLEMENTED)
}

func loadJson(factory *factory.Factory) (*Credential, error) {
	return nil, errors.New(ERR_JSON_FUNCTIONALITY_NOT_IMPLEMENTED)
}
