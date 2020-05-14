package credential

const ERR_USERNAME_OR_PASSWORD_NOT_SET = "sorry you must set a username or password"
const ERR_FACTORY_MUST_BE_INITIALIZED = "sorry the factory must be created correctly before proceeding"
const ERR_KEY_MUST_MATCH_REGEX = "sorry the key must only include numbers, letters and underscores [0-9A-Za-z_]"
const ERR_NOT_INITIALIZED = "sorry something has not been initialized, make sure you use New and do not create a struct manually"
const ERR_CANNOT_SET_USERNAME_WHEN_USING_SECTION = "you cannot set username via this method when using the Section() method"
const ERR_CANNOT_SET_PASSWORD_WHEN_USING_SECTION = "you cannot set password via this method when using the Section() method"
const ERR_CANNOT_REMOVE_USERNAME = "you cannot remove the username from the Credential"
const ERR_CANNOT_REMOVE_PASSWORD = "you cannot remove the username from the Credential"