package global

// Test Labels
//
// These are used as labels when indexing maps within tests.
const TEST_VAR_APPLICATION_NAME = "mtca" // my_test_credential_application
const TEST_VAR_USERNAME_LABEL = "username"
const TEST_VAR_PASSWORD_LABEL = "password"
const TEST_VAR_ATTRIBUTE_NAME_LABEL = "a_test_attribute"
const TEST_VAR_USERNAME_ALTERNATE_LABEL = "access_token"
const TEST_VAR_PASSWORD_ALTERNATE_LABEL = "secret_key"

// Test Environment Labels
//
// These are used as labels when accessing environment variables within tests.
const TEST_VAR_ENVIRONMENT_APPLICATION_NAME = "MTCA"
const TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL = "MTCA_A_TEST_ATTRIBUTE"
const TEST_VAR_ENVIRONMENT_USERNAME_LABEL = "MTCA_USERNAME"
const TEST_VAR_ENVIRONMENT_PASSWORD_LABEL = "MTCA_PASSWORD"
const TEST_VAR_ENVIRONMENT_USERNAME_ALTERNATE_LABEL = "MTCA_ACCESS_TOKEN"
const TEST_VAR_ENVIRONMENT_PASSWORD_ALTERNATE_LABEL = "MTCA_SECRET_KEY"
const TEST_VAR_ENVIRONMENT_BAD_LABEL = "bad key with no underscores"

// Test Environment Labels
//
// These are used as values for variables within tests.
const TEST_VAR_USERNAME = "a_test_username"
const TEST_VAR_PASSWORD = "as=/sle\\sowkjg@!"
const TEST_VAR_BAD_ATTRIBUTE_NAME = "a_ /test_attribute"
const TEST_VAR_ATTRIBUTE_VALUE = "a global attribute value"
const TEST_VAR_USERNAME_ALTERNATE = "another_test_username"
const TEST_VAR_PASSWORD_ALTERNATE = ".YaJ5XAA${hh8^C"

// Variables for Profiles
//
// These are the values that were introduced for building profiles.
const TEST_VAR_FIRST_SECTION_KEY = "first_section"
const TEST_VAR_SECOND_SECTION_KEY = "second_section"
const TEST_VAR_BAD_SECTION_KEY = "bad section key"

const TEST_VAR_DUPLICATE_KEY_LABEL = "duplicate_key"
const TEST_VAR_FIRST_SECTION_UNIQUE_KEY_LABEL = "first_section_key"
const TEST_VAR_SECOND_SECTION_UNIQUE_KEY_LABEL = "second_section_key"
const TEST_VAR_NO_SECTION_UNIQUE_KEY_LABEL = "no_section_key"
const TEST_VAR_BAD_KEY_LABEL = "bad section key label"

const TEST_VAR_DUPLICATE_KEY_VALUE = "duplicate key value"
const TEST_VAR_FIRST_SECTION_UNIQUE_KEY_VALUE = "first section unique value"
const TEST_VAR_SECOND_SECTION_UNIQUE_KEY_VALUE = "second section unique value"
const TEST_VAR_NO_SECTION_UNIQUE_KEY_VALUE = "no section unique value"

const TEST_VAR_FIRST_PROFILE_LABEL = "first_profile"
const TEST_VAR_SECOND_PROFILE_LABEL = "second_profile"
const TEST_VAR_BAD_PROFILE_LABEL = "this is a bad profile name"
