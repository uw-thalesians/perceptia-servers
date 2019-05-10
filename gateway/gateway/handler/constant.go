package handler

import "errors"

// HTTP Header Names.
const (
	HeaderContentType     = "Content-Type"
	HeaderACAllowOrigin   = "Access-Control-Allow-Origin"
	HeaderACAllowMethods  = "Access-Control-Allow-Methods"
	HeaderACAllowHeaders  = "Access-Control-Allow-Headers"
	HeaderACExposeHeaders = "Access-Control-Expose-Headers"
	HeaderACMaxAge        = "Access-Control-Max-Age"
)

// Custom HTTP Header Names
const (
	HeaderPerceptiaUserUuid    = "Perceptia-User-Uuid"
	HeaderPerceptiaSessionUuid = "Perceptia-SessionUuid"
)

// HTTP Content-Type Header Values.
const (
	ContentTypeJSON      = "application/json"
	ContentTypeTextPlain = "text/plain"
)

// HTTP Access-Control Header Values.
const (
	ACAllowOriginAll = "*"
	ACAllowMethods   = "GET, PUT, POST, PATCH, DELETE"
	ACAllowHeaders   = "Content-Type, Authorization"
	ACExposeHeaders  = "Authorization"
	ACMaxAge         = "600"
)

// URL path values.
const SpecificSessionHandlerDeleteUserAlias = "this"

// Handler plain text messages.
const messageSignedOut = "signed out"

// Handler Error Constants.
var (
	errUnexpected = errors.New("an unexpected error has occurred, try again if request did not complete")

	errInvalidCredentials = errors.New("invalid credentials")
	errMethodNotAllowed   = errors.New("method not allowed")

	errMajorVersionNotSupported = errors.New("major version not supported")

	errUserNotFound = errors.New("user not found")
	//errInvalidEmail               = errors.New("invalid email")
	errAccountUserNameUnavailable = errors.New("username unavailable, please select a different user name")

	errActionNotAuthorized = errors.New("action not authorized for the requested resource")
	errUnauthorized        = errors.New("user not authorized, please start a new session")
	errContentTypeNotJson  = errors.New("expected content type was json but application/json content type not set")
)

// request variables
const (
	ReqVarMajorVersion = "majorVersion"
	ReqVarUserUuid     = "userUuid"
	ReqVarSession      = "sessionVar"
)
