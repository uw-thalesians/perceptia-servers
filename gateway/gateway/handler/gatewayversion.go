package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/utility"
)

var ErrPerceptiaApiVersionNotSet = errors.New("perceptia API version not set in request")

var ErrPerceptiaApiVersionNotValidSemVer = errors.New("perceptia API version not a valid sem ver")

var ErrPerceptiaApiVersionNotSupported = errors.New("perceptia API version specified is not supported")

// GatewayVersion represents the current handler in the request/response cycle.
type GatewayVersion struct {
	handler http.Handler
	cx      *Context
}

// NewGatewayVersion constructs a new GatewayVersion struct with the provided handler and Context.
func (cx *Context) NewGatewayVersion(handler http.Handler) http.Handler {
	return &GatewayVersion{handler, cx}
}

// ServeHTTP adds the Perceptia-Api-Version custom header along with the latest version of the gateway api implemented.
func (gv *GatewayVersion) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(HeaderPerceptiaApiVersion, gv.cx.gatewayVersion.String())
	//call the real handler
	gv.handler.ServeHTTP(w, r)
}

// EnsureGatewayVersionSupported represents the current handler in the request/response cycle.
type EnsureGatewayVersionSupported struct {
	handler http.Handler
	cx      *Context
}

// NewEnsureGatewayVersionSupported constructs a new EnsureGatewayVersionSupported struct
// with the provided handler and Context.
func (cx *Context) NewEnsureGatewayVersionSupported(handler http.Handler) http.Handler {
	return &EnsureGatewayVersionSupported{handler, cx}
}

// ServeHTTP checks for the api version header or query param and verifies gateway supports it
func (egv *EnsureGatewayVersionSupported) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqApiVer, errRAP := egv.cx.GetPerceptiaApiVersionValue(r)
	if errRAP != nil {
		if errRAP != ErrPerceptiaApiVersionNotSet {
			if errRAP == ErrPerceptiaApiVersionNotValidSemVer {
				retErr := &Error{
					ClientError: true,
					ServerError: false,
					Message:     fmt.Sprintf("api version specified, but could not be understood, not a valid api version"),
					Context:     "request made to gateway collection",
					Code:        0,
				}
				egv.cx.handleErrorJson(w, r, errRAP, "trying to convert api version header/query param to SemVer", retErr, http.StatusBadRequest)
				return
			} else {
				return
			}
		}

	} else {
		if !egv.cx.ApiVersionSupported(reqApiVer) {
			retErr := &Error{
				ClientError: true,
				ServerError: false,
				Message:     fmt.Sprintf("api version not supported, requested version: %s supported versions: %s", reqApiVer.String(), getSupportedVersionsString(egv.cx.gatewayVersionsSupported)),
				Context:     "request made to gateway collection",
				Code:        0,
			}
			egv.cx.handleErrorJson(w, r, errRAP, "trying to convert api version header/query param to SemVer", retErr, http.StatusBadRequest)
			return
		}
	}

	//call the real handler
	egv.handler.ServeHTTP(w, r)
}

func (cx *Context) GetPerceptiaApiVersionValue(r *http.Request) (apiVer *utility.SemVer, err error) {
	apiVerRaw := r.Header.Get(HeaderPerceptiaApiVersion)
	if apiVerRaw == "" {
		apiVerRawValues, paramSet := (r.URL.Query())[QpApiVersion]
		if !paramSet {
			return nil, ErrPerceptiaApiVersionNotSet
		}
		apiVerRaw := apiVerRawValues[0]
		if apiVerRaw == "" {
			return nil, ErrPerceptiaApiVersionNotSet
		}
	}
	apiVerSVS, errSVS := utility.SemVerFromString(apiVerRaw)
	if errSVS != nil {
		return nil, ErrPerceptiaApiVersionNotValidSemVer
	}
	return apiVerSVS, nil
}

func (cx *Context) ApiVersionSupported(reqApiVersion *utility.SemVer) bool {
	if cx.gatewayVersion.GetMajor() == 0 {
		if cx.gatewayVersion.Compare(reqApiVersion) == 0 {
			return true
		} else {
			if cx.gatewayVersion.GetMinor() == reqApiVersion.GetMinor() {
				if cx.gatewayVersion.GetPatch() >= reqApiVersion.GetPatch() {
					return true
				} else {
					return false
				}
			} else {
				// In major version 0, if minor version don't match, not compatible
				return false
			}
		}
	} else {
		gwVerSel, verSupported := cx.gatewayVersionsSupported[reqApiVersion.GetMajor()]
		if !verSupported {
			return false
		}
		return gwVerSel.Compare(reqApiVersion) >= 0
	}
}

func getSupportedVersionsString(versions map[int]*utility.SemVer) string {
	var versionsString string = ""
	for _, semVer := range versions {
		versionsString = fmt.Sprintf("%s %s", versionsString, semVer.String())
	}
	return versionsString
}
