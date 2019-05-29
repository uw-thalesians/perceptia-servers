package handler

import (
	"net/http"
)

type HealthHandlerContext struct {
	cx                        *Context
	userStoreStatusNotOkay    chan bool
	sessionStoreStatusNotOkay chan bool
}

func (cx *Context) NewHealthHandlerContext(userStoreStatusNotOkay chan bool, sessionStoreStatusNotOkay chan bool) *HealthHandlerContext {
	return &HealthHandlerContext{cx: cx, userStoreStatusNotOkay: userStoreStatusNotOkay, sessionStoreStatusNotOkay: sessionStoreStatusNotOkay}
}

func (hh *HealthHandlerContext) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(HeaderContentType, ContentTypeJSON)
	type healthObj struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}
	gatewayStatus := "ready"
	select {
	case status, ok := <-hh.sessionStoreStatusNotOkay:
		if ok {
			if status {
				gatewayStatus = "not ready"
				hh.sessionStoreStatusNotOkay <- true
			}
		}
	case status, ok := <-hh.userStoreStatusNotOkay:
		if ok {
			if status {
				gatewayStatus = "not ready"
				hh.userStoreStatusNotOkay <- true
			}
		}
	default:
		break
	}

	healthStatus := healthObj{
		Name:   "Perceptia API Health Report",
		Status: gatewayStatus,
	}
	_, _ = hh.cx.respondEncode(w, healthStatus, http.StatusOK)
	return
}
