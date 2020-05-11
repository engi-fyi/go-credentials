package credential

import (
	"errors"
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"github.com/engi-fyi/go-credentials/profile"
	"github.com/engi-fyi/go-credentials/serializer"
)

/*
Load uses LoadFromProfile to create a Credential set to the default profile.
*/
func Load(sourceFactory *factory.Factory) (*Credential, error) {
	return LoadFromProfile(global.DEFAULT_PROFILE_NAME, sourceFactory)
}

/*
Save is responsible for saving the credential at ~/.application_name/credentials in the specified output format
that has been set on the Credentials' Factory object.
*/
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

/*
LoadFromProfile uses Serializer to load an object from the relevant source. The source is determined based on the
OUTPUT_TYPE of the sourceFactory.
*/
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

/*
Serialize retrieves the data values from a Credential object.

This function is designed to be used as part of a serializer.Serialize() call.
*/
func (thisCredential *Credential) Serialize() (string, string, map[string]map[string]string) {
	return thisCredential.Username,
		thisCredential.Password,
		thisCredential.Profile.GetAllAttributes()
}

/*
Deserialize takes all of the values required to build a Credential and Profile and uses those values to do so.

This function is designed to be used as part of a serializer.Deserialize() call.
*/
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
