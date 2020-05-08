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
