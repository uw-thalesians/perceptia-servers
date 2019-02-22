// +build all unit

package utility

import (
	"os"
	"strings"
	"testing"
)

func TestRequireEnv(t *testing.T) {
	var (
		envVarName           = "TEST_VAR_SET"
		envVarNameNotSet     = "TEST_NOT_SET"
		envVarNameValueEmpty = "TEST_VAR_VALUE_EMPTY"
		envVarValue          = "testValueX"
	)
	cases := []struct {
		name                string
		hint                string
		envVarName          string
		expectedEnvVarValue string
		expectError         bool
	}{
		{
			name:                "Basic: env properly set",
			hint:                "Should return the value set for the provided env variable",
			envVarName:          envVarName,
			expectedEnvVarValue: envVarValue,
			expectError:         false,
		},
		{
			name:                "Basic: env set but empty",
			hint:                "Should return error as env value not expected",
			envVarName:          envVarNameValueEmpty,
			expectedEnvVarValue: envVarValue,
			expectError:         true,
		},
		{
			name:                "Basic: env not set",
			hint:                "Should return error as env var not set",
			envVarName:          envVarNameNotSet,
			expectedEnvVarValue: "",
			expectError:         true,
		},
	}
	if err := os.Setenv(envVarName, envVarValue); err != nil {
		t.Errorf("case: N/A: unexpected error in test setup: %s", err)
	}
	if err := os.Setenv(envVarNameValueEmpty, ""); err != nil {
		t.Errorf("case: N/A: unexpected error in test setup: %s", err)
	}

	for _, c := range cases {
		returnedValue, errRE := RequireEnv(c.envVarName)
		if errRE != nil && !c.expectError {
			t.Errorf("case: %s: error not expected but got %s\nHINT: %s", c.name, errRE, c.hint)
		} else if c.expectError {
			if errRE == nil {
				t.Errorf("case: %s: expected error but got %t\nHINT: %s", c.name, errRE, c.hint)
			}
		} else if strings.Compare(returnedValue, c.expectedEnvVarValue) != 0 {
			t.Errorf("case: %s: expected set value: %s, but got: %s\nHINT: %s", c.name, c.expectedEnvVarValue, returnedValue, c.hint)
		}
	}
}
