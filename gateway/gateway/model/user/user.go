package user

import (
	"errors"
	"net/mail"
	"strings"
)

// Validation constants
const (
	ValidPasswordMinLength    = 8
	ValidUserNameMinLength    = 3
	ValidUserNameMaxLength    = 255
	ValidFullNameMaxLength    = 255
	ValidDisplayNameMaxLength = 255
)

// Custom Error types
var (
	// ErrEmailInvalid used when the provided email is not valid.
	ErrEmailInvalid = errors.New("invalid email")

	// ErrPasswordLengthLessThanMin used when the provided password does not meet minimum length requirements.
	ErrPasswordLengthLessThanMin = errors.New("password must be at least " +
		string(ValidPasswordMinLength) + " characters long")

	// ErrUserNameLengthLessThanMin used when the provided userName is not long enough.
	ErrUserNameLengthLessThanMin = errors.New("username must be at least " +
		string(ValidUserNameMinLength) + " characters long")

	// ErrUserNameLengthGreaterThanMax used when the provided userName is too long.
	ErrUserNameLengthGreaterThanMax = errors.New("username must be less than " +
		string(ValidUserNameMaxLength) + " characters long")

	// ErrUserNameHasSpace used when the provided userName has spaces.
	ErrUserNameHasSpace = errors.New("username must not have any spaces")

	// ErrFullNameLengthGreaterThanMax used when the provided fullName is too long.
	ErrFullNameLengthGreaterThanMax = errors.New("full name must be less than " +
		string(ValidFullNameMaxLength) + " characters long")

	// ErrDisplayNameLengthGreaterThanMax used when the provided displayName is too long.
	ErrDisplayNameLengthGreaterThanMax = errors.New("display name must be less than " +
		string(ValidDisplayNameMaxLength) + " characters long")

	// ErrHashNotFromPassword used when the provided password was not the password used to create the user's PassHash.
	ErrHashNotFromPassword = errors.New("the provided password did not create this user's PassHash")
)

// User represents a user account in the database.
type User struct {
	UUID        string `json:"uuid"`
	Email       string `json:"-"` //never JSON encoded/decoded
	UserName    string `json:"userName"`
	FullName    string `json:"fullName"`
	DisplayName string `json:"displayName"`
}

// SignInCredentials represents user sign-in credentials.
type SignInCredentials struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// Credential represents the user's stored secret
type Credential struct {
	EncodedPassword string `json:"-"` // never JSON encoded/decoded
}

// NewUser represents a new user signing up for an account.
type NewUser struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	UserName    string `json:"userName"`
	FullName    string `json:"fullName"`
	DisplayName string `json:"displayName"`
}

// Updates represents allowed updates to a user profile.
type Updates struct {
	FullName    string `json:"fullName"`
	DisplayName string `json:"displayName"`
}

// Validate validates the new user and returns an error if any of the validation rules fail, or nil if its valid.
//
// Validation rules: (Only one error will be returned if multiple validation errors are present;
// fail order is not guaranteed):
// - Email field must be a valid email address.
// - Password must be at least 6 characters.
// - Password and PasswordConf must match.
// - UserName must be non-zero length and may not contain spaces.
func (nu *NewUser) Validate() error {
	if err := validateEmail(nu.Email); err != nil {
		return err
	}
	if err := validateUserName(nu.UserName); err != nil {
		return err
	}
	if err := validateFullName(nu.FullName); err != nil {
		return err
	}
	if err := validateDisplayName(nu.DisplayName); err != nil {
		return err
	}
	return nil
}

// ToUser converts the NewUser to a User, setting the PhotoURL and PassHash fields appropriately.
func (nu *NewUser) ToUser() (*User, error) {
	errV := nu.Validate()
	if errV != nil {
		return nil, errV
	}
	usr := &User{
		Email:       strings.TrimSpace(nu.Email),
		UserName:    nu.UserName,
		FullName:    strings.TrimSpace(nu.FullName),
		DisplayName: strings.TrimSpace(nu.DisplayName),
	}

	errSP := usr.SetPassword(nu.Password)
	if errSP != nil {
		return nil, errSP
	}

	return usr, nil
}

// SetPassword hashes the password and stores it in the PassHash field.
func (u *User) SetPassword(password string) error {
	// TODO
	return nil
}

// Authenticate compares the plaintext password against the stored hash and returns an error if
// they don't match, or nil if they do.
func (u *User) Authenticate(password string) error {
	// TODO
	return nil
}

// ApplyUpdates applies the updates to the user. An error is returned if the updates are invalid.
func (u *User) ApplyUpdates(updates *Updates) error {
	u.FullName = strings.TrimSpace(updates.FullName)
	u.DisplayName = strings.TrimSpace(updates.DisplayName)
	return nil
}

// ValidateCredentialsEmail ensures that the supplied email is a valid email.
// A valid email is one that would pass the Validate function for a NewUser.
// Will return nil if the email is valid as defined above.
func (c *Credentials) ValidateCredentialsEmail() error {
	return validateEmail(c.Email)
}

// validateEmail validates the provided email.
// If valid, returns nil, otherwise an error.
func validateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return ErrEmailInvalid
	}
	return nil
}

// validatePassword validates the provided password.
// If valid, returns nil, otherwise an error.
func validatePassword(password string) error {
	if len([]rune(password)) < ValidPasswordMinLength {
		return ErrPasswordLengthLessThanMin
	}
	return nil
}

// validateUserName validates the provided userName.
// If valid, returns nil, otherwise an error.
// (If multiple validation errors occur only one error will be returned; order of validation is not guarantied)
func validateUserName(userName string) error {
	if strings.Contains(userName, " ") {
		return ErrUserNameHasSpace
	}
	lenUserName := len([]rune(userName))
	if lenUserName < ValidUserNameMinLength {
		return ErrUserNameLengthLessThanMin
	} else if lenUserName > ValidUserNameMaxLength {
		return ErrUserNameLengthGreaterThanMax
	}
	return nil
}

// validateFullName validates the provided fullName.
// If valid, returns nil, otherwise an error.
// (If multiple validation errors occur only one error will be returned; order of validation is not guarantied)
func validateFullName(fullName string) error {
	if len([]rune(fullName)) > ValidFullNameMaxLength {
		return ErrFullNameLengthGreaterThanMax
	}
	return nil
}

// validateDisplayName validates the provided displayName.
// If valid, returns nil, otherwise an error.
// (If multiple validation errors occur only one error will be returned; order of validation is not guarantied)
func validateDisplayName(displayName string) error {

	if len([]rune(displayName)) > ValidFullNameMaxLength {
		return ErrDisplayNameLengthGreaterThanMax
	}
	return nil
}
