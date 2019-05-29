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
	HeaderAuthorization   = "Authorization"
	HeaderLocation        = "Location"
	HeaderCacheControl    = "Cache-Control"
	HeaderPragma          = "Pragma"
	HeaderContentLength   = "Content-Length"
	HeaderWWWAuthenticate = "WWW-Authenticate"
	// Custom HTTP Header Names
	HeaderPerceptiaUserUuid    = "Perceptia-User-Uuid"
	HeaderPerceptiaSessionUuid = "Perceptia-Session-Uuid"
	HeaderPerceptiaApiVersion  = "Perceptia-Api-Version"
)

const (
	// HTTP Content-Type Header Values.
	ContentTypeJSON      = "application/json"
	ContentTypeTextPlain = "text/plain"
	// HTTP Access-Control Header Values.
	ACAllowOriginAll = "*"
	ACAllowMethods   = "GET, PUT, POST, PATCH, DELETE"
	ACAllowHeaders   = HeaderContentType + ", " + HeaderAuthorization + ", " + HeaderPerceptiaApiVersion
	ACExposeHeaders  = HeaderAuthorization + ", " + HeaderPerceptiaApiVersion + ", " + HeaderLocation + ", " +
		HeaderCacheControl + ", " + HeaderPragma + ", " + HeaderContentLength + ", " + HeaderWWWAuthenticate
	ACMaxAge = "600"
	// HTTP Cache and Pragma Header Values
	CacheControlNoStore = "no-store"
	PragmaNoCache       = "no-cache"
	// AuthorizationSchemeValues
	AuthorizationBearer = "Bearer"
	// rfc6750#section-3
	WWWAuthenticateBearerRealm            = "Bearer realm=\"/api/\""
	WWWAuthenticateErrorInvalidToken      = "error=\"invalid_token\""
	WWWAuthenticateErrorInvalidRequest    = "error=\"invalid_request\""
	WWWAuthenticateErrorInsufficientScope = "error=\"insufficient_scope\""
)

// Query Parameters
const (
	QpApiVersion = "apiVersion"
)

// URL path values.
const SpecificSessionHandlerDeleteCurrentSessionAlias = "this"

// Handler Error Constants.
var (
	errUnexpected = errors.New("an unexpected error has occurred, try again if request did not complete")

	errInvalidCredentials = errors.New("invalid credentials")
	errMethodNotAllowed   = errors.New("method not allowed")

	errMajorVersionNotSupported = errors.New("major version not supported")

	errUserNotFound = errors.New("user not found")
	//errInvalidEmail               = errors.New("invalid email")
	errAccountUserNameUnavailable = errors.New("username unavailable, please select a different user name")
	errSessionNotFound            = errors.New("session not found")
	errUserNotInSession           = errors.New("not in a session")

	errActionNotAuthorized = errors.New("action not authorized for the requested resource")
	errUnauthorized        = errors.New("user not authorized, please start a new session")
	errContentTypeNotJson  = errors.New("expected content type was json but application/json content type not set")
	errDecodingJson        = errors.New("issue decoding request body into json object")
)

// Gmux request variables
const (
	ReqVarMajorVersion = "majorVersion"
	ReqVarUserUuid     = "userUuid"
	ReqVarSession      = "sessionVar"
)
