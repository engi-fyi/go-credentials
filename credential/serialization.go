package credential

import (
	"errors"
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"github.com/engi-fyi/go-credentials/profile"
	"github.com/engi-fyi/go-credentials/serializer"
)

/*
Load is the default method of loading existing credentials. First it will attempt to load credentials from the
environment, as they take precedence over file-based credentials. If that is not successful, it will attempt to
load credentials from the credentials file using the appropriate format set in the sourceFactory. If the file doesn't
exist or is in an unexpected format, an error will be thrown.
*/
// TODO(7): Implement json format.
func Load(sourceFactory *factory.Factory) (*Credential, error) {
	return LoadFromProfile(global.DEFAULT_PROFILE_NAME, sourceFactory)
}

func (thisCredential *Credential) Save() error {
	if !thisCredential.Factory.Initialized || !thisCredential.Initialized {
		return errors.New(ERR_NOT_INITIALIZED)
	}

	if !thisCredential.Profile.Initialized {
		return errors.New(profile.ERR_PROFILE_NOT_INITIALIZED)
	}

		mySerializer := serializer.New(thisCredential.Factory, thisCredential.Profile.Name)
	return mySerializer.Serialize(thisCredential.Serialize())
}

func LoadFromProfile(profileName string, sourceFactory *factory.Factory) (*Credential, error) {
	if !sourceFactory.Initialized {
		return nil, errors.New(ERR_FACTORY_MUST_BE_INITIALIZED)
	}

	mySerializer := serializer.New(sourceFactory, profileName)
	username, password, attributes, deErr := mySerializer.Deserialize()

	if deErr != nil {
		return nil, deErr
	}

	return Deserialize(sourceFactory, profileName, username, password, attributes)
}

func (thisCredential *Credential) Serialize() (string, string, map[string]map[string]string) {
	return thisCredential.Username,
		thisCredential.Password,
		thisCredential.Profile.GetAllAttributes()
}

func Deserialize(sourceFactory *factory.Factory, profileName string, username string, password string, attributes map[string]map[string]string) (*Credential, error) {
	myCredential, credErr := New(sourceFactory, username, password)
	myProfile, profileErr := profile.New(profileName, sourceFactory)

	if credErr != nil {
		return nil, credErr
	}

	if profileErr != nil {
		return nil, profileErr
	}

	for section := range attributes {
		for key, value := range attributes[section] {
			setErr := myProfile.SetAttribute(section, key, value)

			if setErr != nil {
				return nil, setErr
			}
		}
	}

	myCredential.Profile = myProfile
	return myCredential, nil
}