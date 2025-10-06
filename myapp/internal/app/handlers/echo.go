package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/RemeshevskiyValeriy/myapp/utils"
)


func Echo(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)

	if r.Method != http.MethodPost {
		utils.WriteErr(w, http.StatusMethodNotAllowed, "only POST allowed")
		return
	}

	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteErr(w, http.StatusBadRequest, "invalid JSON" + err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(body)
}