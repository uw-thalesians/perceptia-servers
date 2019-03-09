// +build all unit

package user

import "testing"

// TODO: Write tests for user

// TestNewUser_ValidateNewUser is a unit test ensuring ValidateNewUser
// check all fields of the struct.
//
// This test does not test all possible inputs; more exhaustive tests will
// be run on the individual validation functions used by this method.
func TestNewUser_ValidateNewUser(t *testing.T) {
	tests := []struct {
		name        string
		nu          *NewUser
		detail      string
		expectError bool
	}{
		{"Basic Test",
			&NewUser{
				"test",
				"Test Tester",
				"Test",
				"$argon2id$v=19$m=65536,t=1," +
					"p=2$9QWBW808Dq62aOfojSwBYg$muj80qzHXdB4DW06zgHbljhTLGrfFJ6qq8hDiK" +
					"2esNTUBH0zlnBjGoDJ4S5glnKVCs9CD15cajhzkVxjgMwgHZdZ0JgZvUDcX9tJf8p7a" +
					"TCOI8vN0vVVbZFbgUya6CFiqPH8cULMVyDJLK84hzqgWlryxt4p5eZjladgHyeKPlw",
			},
			"NewUser Valid. Should return nil",
			false,
		},
		{"Missing Username",
			&NewUser{
				"",
				"Test Tester",
				"Test",
				"$argon2id$v=19$m=65536,t=1," +
					"p=2$9QWBW808Dq62aOfojSwBYg$muj80qzHXdB4DW06zgHbljhTLGrfFJ6qq8hDiK" +
					"2esNTUBH0zlnBjGoDJ4S5glnKVCs9CD15cajhzkVxjgMwgHZdZ0JgZvUDcX9tJf8p7a" +
					"TCOI8vN0vVVbZFbgUya6CFiqPH8cULMVyDJLK84hzqgWlryxt4p5eZjladgHyeKPlw",
			},
			"NewUser Invalid. Should return an error",
			true,
		},
		{"Missing FullName",
			&NewUser{
				"test",
				"",
				"Test",
				"$argon2id$v=19$m=65536,t=1," +
					"p=2$9QWBW808Dq62aOfojSwBYg$muj80qzHXdB4DW06zgHbljhTLGrfFJ6qq8hDiK" +
					"2esNTUBH0zlnBjGoDJ4S5glnKVCs9CD15cajhzkVxjgMwgHZdZ0JgZvUDcX9tJf8p7a" +
					"TCOI8vN0vVVbZFbgUya6CFiqPH8cULMVyDJLK84hzqgWlryxt4p5eZjladgHyeKPlw",
			},
			"NewUser Valid. Should return nil",
			false,
		},
		{"Missing DisplayName",
			&NewUser{
				"test",
				"Test Tester",
				"",
				"$argon2id$v=19$m=65536,t=1," +
					"p=2$9QWBW808Dq62aOfojSwBYg$muj80qzHXdB4DW06zgHbljhTLGrfFJ6qq8hDiK" +
					"2esNTUBH0zlnBjGoDJ4S5glnKVCs9CD15cajhzkVxjgMwgHZdZ0JgZvUDcX9tJf8p7a" +
					"TCOI8vN0vVVbZFbgUya6CFiqPH8cULMVyDJLK84hzqgWlryxt4p5eZjladgHyeKPlw",
			},
			"NewUser Valid. Should return nil",
			false,
		},
		{"Missing EncodedHash",
			&NewUser{
				"test",
				"Test Tester",
				"Test",
				"",
			},
			"NewUser Invalid. Should return an error",
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errVNU := test.nu.ValidateNewUser()
			if test.expectError && errVNU == nil {
				t.Errorf("An error was expected, but did not occur\nSee detail:\n\t%s", test.detail)
			} else if !test.expectError && errVNU != nil {
				t.Errorf("An error was not expected, but one occured.\n\t"+
					"Error: %s\n\tDetail: %s", errVNU, test.detail)
			}
		})
	}
}

// TestNewUser_PrepNewUser is a unit test ensuring that PrepNewUser prepares all
// fields that should be prepared.
//
// This test does not test all possible inputs; more exhaustive tests will
// be run on the individual prep functions used by this method.
func TestNewUser_PrepNewUser(t *testing.T) {
	tests := []struct {
		name       string
		nu         *NewUser
		expectedNu *NewUser
		detail     string
	}{
		{"Basic Test",
			&NewUser{
				"test",
				"Test Tester",
				"Test",
				"",
			},
			&NewUser{
				"test",
				"Test Tester",
				"Test",
				"",
			},
			"Both structs the same, all fields should equal.",
		},
		{"Missing Username",
			&NewUser{
				"",
				"Test Tester",
				"Test",
				"",
			},
			&NewUser{
				"",
				"Test Tester",
				"Test",
				"",
			},
			"Both structs the same, all fields should equal.",
		},
		{"All Fields Missing",
			&NewUser{
				"",
				"",
				"",
				"",
			},
			&NewUser{
				"",
				"",
				"",
				"",
			},
			"Both structs the same, all fields should equal.",
		},
		{"Extra Spaces around fields",
			&NewUser{
				" test  ",
				" Test Tester   ",
				" Test ",
				"",
			},
			&NewUser{
				"test",
				"Test Tester",
				"Test",
				"",
			},
			"Both structs the same, all fields should equal.",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.nu.PrepNewUser()
			if test.nu.Username != test.expectedNu.Username {
				t.Errorf("Usernames do not match. Got: %s, Expected: %s",
					test.nu.Username, test.expectedNu.Username)
			}
			if test.nu.FullName != test.expectedNu.FullName {
				t.Errorf("FullNames do not match. Got: %s, Expected: %s",
					test.nu.FullName, test.expectedNu.FullName)
			}
			if test.nu.DisplayName != test.expectedNu.DisplayName {
				t.Errorf("DisplayNames do not match. Got: %s, Expected: %s",
					test.nu.DisplayName, test.expectedNu.DisplayName)
			}
		})
	}
}
