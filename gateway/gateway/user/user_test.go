// +build all unit

package user

import "testing"

// TODO: Write tests for user

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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errVNU := test.nu.ValidateNewUser()
			if test.expectError && errVNU == nil {
				t.Errorf("An error was expected, but did not occur\nSee detail:\n\t%s", test.detail)
			}
		})
	}
}
