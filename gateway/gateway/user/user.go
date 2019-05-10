// Package user provides user structs and database implementations for storing user data
package user

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	uuid "github.com/satori/go.uuid"

	"golang.org/x/crypto/argon2"
)

// Validation constants
const (
	ValidPasswordMinLength    = 8
	ValidPasswordMaxLength    = 500
	ValidUsernameMinLength    = 3
	ValidUsernameMaxLength    = 255
	ValidFullNameMaxLength    = 255
	ValidDisplayNameMaxLength = 255
)

const InvalidEncodedPasswordHash = ""

const InvalidEmail = ""

// Custom Error types
var (
	// ErrInvalidEmail used when the provided email is not valid.
	ErrInvalidEmail = errors.New("invalid email")

	// ErrPasswordLengthLessThanMin used when the provided password does not
	// meet minimum length requirements.
	ErrPasswordLengthLessThanMin = fmt.Errorf("password must be at least %d characters long",
		ValidPasswordMinLength)

	// ErrPasswordLengthGreaterThanMax used when the provided password is too long.
	ErrPasswordLengthGreaterThanMax = errors.New("password must be no more than " +
		string(ValidPasswordMaxLength) + " characters long")

	// ErrUsernameLengthLessThanMin used when the provided username is not long enough.
	ErrUsernameLengthLessThanMin = errors.New("username must be at least " +
		string(ValidUsernameMinLength) + " characters long")

	// ErrUsernameLengthGreaterThanMax used when the provided username is too long.
	ErrUsernameLengthGreaterThanMax = errors.New("username must be no more than " +
		string(ValidUsernameMaxLength) + " characters long")

	// ErrUserNameHasSpace used when the provided username has spaces.
	ErrUserNameHasSpace = errors.New("username must not have any spaces")

	// ErrFullNameLengthGreaterThanMax used when the provided fullName is too long.
	ErrFullNameLengthGreaterThanMax = errors.New("full name must be no more than " +
		string(ValidFullNameMaxLength) + " characters long")

	// ErrDisplayNameLengthGreaterThanMax used when the provided displayName is too long.
	ErrDisplayNameLengthGreaterThanMax = errors.New("display name must be no more than " +
		string(ValidDisplayNameMaxLength) + " characters long")

	// ErrHashNotFromPassword used when the provided password was not
	// the password used to create the user's EncodedHash.
	ErrHashNotFromPassword = errors.New("the provided password is not the current password")

	// ErrInvalidCredentials is used when the provided login credentials are invalid
	ErrInvalidCredentials = errors.New("the provided username or password are invalid")

	// ErrInvalidHash indicates that the hash was not encoded correctly
	ErrInvalidHash = errors.New("the encoded hash is not in the correct format")

	// ErrIncompatibleVersion indicates that the hash was created with
	// an incompatible version of argon2
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

// User represents the standard struct for storing basic user information.
type User struct {
	Uuid        uuid.UUID `json:"uuid"`
	Username    string    `json:"username"`
	FullName    string    `json:"fullName"`
	DisplayName string    `json:"displayName"`
}

// SignInCredentials represents user sign-in credentials.
type SignInCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewUser represents a new user signing up for an account.
type NewUser struct {
	Username    string `json:"username"`
	FullName    string `json:"fullName"`
	DisplayName string `json:"displayName"`
	EncodedHash string `json:"encodedHash"`
}

type Updates struct {
	FullName    string `json:"fullName,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

// argon2Params represents the parameters to the Argon2 password hashing algorithm.
type argon2Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// specificArgon2Params are the specific parameters used to create the argon2 hash of the password.
var specificArgon2Params = &argon2Params{
	memory:      64 * 1024,
	iterations:  1,
	parallelism: 2,
	saltLength:  16,
	keyLength:   128,
}

// ValidateNewUser validates the new user and returns an error if any of the validation rules fail,
// or nil if it's valid.
//
// Validation rules: (Only one error will be returned if multiple validation errors are present;
// fail order is not guaranteed):
//
// - Username must be non-zero length and may not contain spaces.
// - FullName must be less than the maximum length for the field.
// - DisplayName must be less than the maximum length for the field.
// - EncodedHash must be set and a valid EncodedHash.
func (nu *NewUser) ValidateNewUser() error {
	if err := ValidateUsername(nu.Username); err != nil {
		return err
	}
	if err := ValidateFullName(nu.FullName); err != nil {
		return err
	}
	if err := ValidateDisplayName(nu.DisplayName); err != nil {
		return err
	}
	if err := ValidateEncodedHash(nu.EncodedHash); err != nil {
		return err
	}
	return nil
}

// PrepNewUser prepares a NewUser struct to be added to the database.
func (nu *NewUser) PrepNewUser() {
	nu.Username = PrepUsername(nu.Username)
	nu.FullName = PrepFullName(nu.FullName)
	nu.DisplayName = PrepDisplayName(nu.DisplayName)
}

func (u *Updates) ValidateUpdates() error {
	if err := ValidateDisplayName(u.DisplayName); err != nil {
		return err
	}
	if err := ValidateFullName(u.FullName); err != nil {
		return err
	}
	return nil
}

func (u *Updates) PrepUpdates() {
	u.FullName = PrepFullName(u.FullName)
	u.DisplayName = PrepDisplayName(u.DisplayName)
}

// PrepFullName prepares the provided string to be used.
func PrepFullName(fullName string) string {
	return strings.TrimSpace(fullName)
}

// PrepUsername prepares the provided string to be used.
func PrepUsername(username string) string {
	return strings.TrimSpace(username)
}

// PrepDisplayName prepares the provided string to be used.
func PrepDisplayName(displayName string) string {
	return strings.TrimSpace(displayName)
}

// CleanEmail returns a valid email address extracted from the provided address,
// or an invalid email and an error if unable to extract the email from the address.
func CleanEmail(address string) (userEmail string, error error) {
	addr, errPA := mail.ParseAddress(address)
	if errPA != nil {
		return InvalidEmail, ErrInvalidEmail
	}
	return addr.Address, nil
}

// CreateEncodedHash takes the provided password and hashes it.
// If the password is invalid or an error occurs while creating the encoded hash,
// the error will be returned along with an invalid encoded password hash.
func CreateEncodedHash(password string) (string, error) {
	errVP := ValidatePassword(password)
	if errVP != nil {
		return InvalidEncodedPasswordHash, errVP
	}
	encodedHash, errGFP := generateFromPassword(password, specificArgon2Params)
	if errGFP != nil {
		return InvalidEncodedPasswordHash, errGFP
	}
	return encodedHash, nil
}

// Authenticate compares the plaintext password against the encoded hash.
// If the password matches with the hashed password true is returned and a nil error.
// If the passwords don't match false is returned, along with ErrHashNotFromPassword
func Authenticate(password, encodedHash string) (bool, error) {
	if bl, _ := comparePasswordAndHash(password, encodedHash); bl == false {
		return false, ErrHashNotFromPassword
	}
	return true, nil
}

// ValidateSignInCredentials ensures that the supplied username is a valid username.
// A valid username is one that would pass the Validate function for a NewUser.
// Will return nil if the username is valid as defined above.
func (c *SignInCredentials) ValidateSignInCredentials() error {
	if ValidateUsername(c.Username) != nil || ValidatePassword(c.Password) != nil {
		return ErrInvalidCredentials
	}
	return ValidateUsername(c.Username)
}

// ValidateEmail validates the provided email.
// If valid, returns nil, otherwise an error.
func ValidateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidEmail
	}
	return nil
}

// ValidatePassword validates the provided password.
// If valid, returns nil, otherwise an error.
func ValidatePassword(password string) error {
	lenPass := len([]rune(password))
	if lenPass < ValidPasswordMinLength {
		return ErrPasswordLengthLessThanMin
	} else if lenPass > ValidPasswordMaxLength {
		return ErrPasswordLengthGreaterThanMax
	}
	return nil
}

// ValidateEncodedHash validates the provided encoded hash.
// If valid, returns nil, otherwise an error.
func ValidateEncodedHash(encodedHash string) error {
	if _, _, _, err := decodeHash(encodedHash); err != nil {
		return ErrInvalidHash
	}
	return nil
}

// ValidateUsername validates the provided username.
// If valid, returns nil, otherwise an error.
// (If multiple validation errors occur only one error will be returned; order of validation is not guarantied)
func ValidateUsername(username string) error {
	if strings.Contains(username, " ") {
		return ErrUserNameHasSpace
	}
	lenUserName := len([]rune(username))
	if lenUserName < ValidUsernameMinLength {
		return ErrUsernameLengthLessThanMin
	} else if lenUserName > ValidUsernameMaxLength {
		return ErrUsernameLengthGreaterThanMax
	}
	return nil
}

// ValidateFullName validates the provided fullName.
// If valid, returns nil, otherwise an error.
// (If multiple validation errors occur only one error will be returned; order of validation is not guarantied)
func ValidateFullName(fullName string) error {
	if len([]rune(fullName)) > ValidFullNameMaxLength {
		return ErrFullNameLengthGreaterThanMax
	}
	return nil
}

// ValidateDisplayName validates the provided displayName.
// If valid, returns nil, otherwise an error.
// (If multiple validation errors occur only one error will be returned; order of validation is not guarantied)
func ValidateDisplayName(displayName string) error {
	if len([]rune(displayName)) > ValidFullNameMaxLength {
		return ErrDisplayNameLengthGreaterThanMax
	}
	return nil
}

// generateFromPassword generates an encoded hash of the provided password
// using the provided Argon2id parameters `p`.
//
// If successful, will return a valid encodedHash string and a nil error.
// If an error occurs while generating the encoded hash,
// will return an InvalidEncodedPasswordHash and the error that
// occurred.
//
// Attribution: https://gist.github.com/alexedwards/34277fae0f48abe36822b375f0f6a621
func generateFromPassword(password string, p *argon2Params) (encodedHash string, err error) {
	salt, err := generateRandomBytes(p.saltLength)
	if err != nil {
		return InvalidEncodedPasswordHash, err
	}

	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// generateRandomBytes generates a byte slice of randomly generated data
// from a cryptographically secure source of randomness.
//
// Attribution: https://gist.github.com/alexedwards/34277fae0f48abe36822b375f0f6a621
func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// comparePasswordAndHash compares the provided password and encoded hash.
//
// If the provided password was not the password that created the provided hash,
// the function will return false. Otherwise the function will return true.
// Any errors that occur will also be returned along with false.
//
// Attribution: https://gist.github.com/alexedwards/34277fae0f48abe36822b375f0f6a621
func comparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

// decodeHash extracts the components of the encoded hash and returns them.
//
// Attribution: https://gist.github.com/alexedwards/34277fae0f48abe36822b375f0f6a621
func decodeHash(encodedHash string) (p *argon2Params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &argon2Params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
