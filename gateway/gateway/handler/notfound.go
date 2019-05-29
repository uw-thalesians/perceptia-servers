package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (cx *Context) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	retErr := &Error{
		ClientError: true,
		ServerError: false,
		Message:     mux.ErrNotFound.Error(),
		Context:     fmt.Sprintf("method=%s path=%s", r.Method, r.URL.Path),
		Code:        0,
	}
	cx.handleErrorJson(w, r, nil, "requested resource not found",
		retErr, http.StatusNotFound)
	return
}
