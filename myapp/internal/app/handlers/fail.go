package handlers

import (
	"net/http"

	"github.com/RemeshevskiyValeriy/myapp/utils"
)

func Fail(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	utils.WriteErr(w, http.StatusBadRequest, "bad_request_example")
}
