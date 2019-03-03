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
	ValidUserNameMinLength    = 3
	ValidUserNameMaxLength    = 255
	ValidFullNameMaxLength    = 255
	ValidDisplayNameMaxLength = 255
)

const InvalidEncodedPasswordHash = ""

// Custom Error types
var (
	// ErrInvalidEmail used when the provided email is not valid.
	ErrInvalidEmail = errors.New("invalid email")

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

	// ErrHashNotFromPassword used when the provided password was not the password used to create the user's EncodedHash.
	ErrHashNotFromPassword = errors.New("the provided password is not the users current password")

	// ErrInvalidCredentials is used when the provided login credentials are invalid
	ErrInvalidCredentials = errors.New("the provided username or password are invalid")

	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
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

// Email represents the address part of a user's email
type Email struct {
	Address string `json:"email"`
}

// UpdateFullName represents an update to the user's full name
type UpdateFullName struct {
	FullName string `json:"fullName"`
}

// UpdateDisplayName represents an update to the user's display name
type UpdateDisplayName struct {
	DisplayName string `json:"displayName"`
}

// argon2Params represents the parameters to the Argon2 password hashing algorithm
type argon2Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var specificArgon2Params = &argon2Params{
	memory:      64 * 1024,
	iterations:  1,
	parallelism: 2,
	saltLength:  16,
	keyLength:   128,
}

// Validate validates the new user and returns an error if any of the validation rules fail, or nil if its valid.
//
// Validation rules: (Only one error will be returned if multiple validation errors are present;
// fail order is not guaranteed):
// - Username must be non-zero length and may not contain spaces.
// - FullName must be less than the maximum length for the field.
// - DisplayName must be less than the maximum length for the field.
// - EncodedHash must be set and a valid EncodedHash.
func (nu *NewUser) Validate() error {
	if err := validateUsername(nu.Username); err != nil {
		return err
	}
	if err := validateFullName(nu.FullName); err != nil {
		return err
	}
	if err := validateDisplayName(nu.DisplayName); err != nil {
		return err
	}
	if err := validateEncodedHash(nu.EncodedHash); err != nil {
		return err
	}
	return nil
}

// PrepNewUser prepares a NewUser struct to be added to the database.
func (nu *NewUser) PrepNewUser() {
	nu.Username = strings.TrimSpace(nu.Username)
	nu.FullName = strings.TrimSpace(nu.FullName)
	nu.DisplayName = strings.TrimSpace(nu.DisplayName)
}

// NewEmail creates a Email object using the provided email address.
// An error is returned if there is an issue with the address.
func NewEmail(address string) (userEmail *Email, error error) {
	addr, errPA := mail.ParseAddress(address)
	if errPA != nil {
		return nil, ErrInvalidEmail
	}
	return &Email{Address: addr.Address}, nil
}

// CreateEncodedHash takes the provided password and hashes it.
// If the password is invalid or an error occurs while creating the encoded hash,
// the error will be returned along with an invalid encoded password hash.
func CreateEncodedHash(password string) (string, error) {
	errVP := validatePassword(password)
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

// ValidateSignInCredentails ensures that the supplied username is a valid username.
// A valid username is one that would pass the Validate function for a NewUser.
// Will return nil if the username is valid as defined above.
func (c *SignInCredentials) ValidateSignInCredentials() error {
	if validateUsername(c.Username) != nil || validatePassword(c.Password) != nil {
		return ErrInvalidCredentials
	}
	return validateUsername(c.Username)
}

// validateEmail validates the provided email.
// If valid, returns nil, otherwise an error.
func validateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidEmail
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

// validateEncodedHash validates the provided encoded hash.
// If valid, returns nil, otherwise an error.
func validateEncodedHash(encodedHash string) error {
	if _, _, _, err := decodeHash(encodedHash); err != nil {
		return ErrInvalidHash
	}
	return nil
}

// validateUserName validates the provided userName.
// If valid, returns nil, otherwise an error.
// (If multiple validation errors occur only one error will be returned; order of validation is not guarantied)
func validateUsername(userName string) error {
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

// generateFromPassword generates an encoded hash of the provided password using the provided Argon2id parameters `p`.
// If successful, will return a valid encodedHash string and a nil error.
// If an error occurs while generating the encoded hash, will return an InvalidEncodedPasswordHash and the error that
// occurred.
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

// generateRandomBytes generates a byte slice of randomly generated data from a cyrptographically secure source of
// randomness.
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
// If the provided password was not the password that created the provided hash,
// the function will return false. Otherwise the function will return true.
// Any errors that occur will also be returned along with false.
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