package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RemeshevskiyValeriy/myapp/utils"
)


func Time(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{
		"time": time.Now().UTC().Format(time.RFC3339),
	})
}