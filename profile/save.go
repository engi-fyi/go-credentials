package profile

import "gopkg.in/ini.v1"

func (thisProfile *Profile) Save() error {
	myFile := ini.Empty()

	for key, value := range thisProfile.attributes {
		mySection, sectionErr := myFile.NewSection(key)

		if sectionErr != nil {
			return sectionErr
		}

		for subKey, subValue := range value {
			_, keyErr := mySection.NewKey(subKey, subValue)

			if keyErr != nil {
				return keyErr
			}
		}
	}

	saveErr := myFile.SaveTo(thisProfile.ConfigFileLocation)

	if saveErr != nil {
		return saveErr
	}

	return nil
}